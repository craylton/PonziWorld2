import type { Asset } from './Asset';
import './AssetList.css';
import ChevronIcon from '../ChevronIcon';
import { useState, useEffect, useCallback } from 'react';
import InvestedAssetSummary from './InvestedAssetSummary';
import UninvestedAssetSummary from './UninvestedAssetSummary';

interface AssetListProps {
    title: string;
    onLoad: () => Promise<Asset[]>;
    isExpandedByDefault: boolean;
}

export default function AssetList({ title, onLoad, isExpandedByDefault }: AssetListProps) {
    const [allAssets, setAllAssets] = useState<Asset[]>([]);
    const [isExpanded, setIsExpanded] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [hasLoadedAssetTypes, setHasLoadedAssetTypes] = useState(false);

    const handleToggleAssets = useCallback(async () => {
        if (isExpanded) {
            setIsExpanded(false);
            return;
        }

        if (hasLoadedAssetTypes) {
            setIsExpanded(true);
            return;
        }

        setIsLoading(true);
        try {
            setAllAssets(await onLoad());
            setHasLoadedAssetTypes(true);
            setIsExpanded(true);
        } catch (error) {
            console.error('Error loading asset types:', error);
        } finally {
            setIsLoading(false);
        }
    }, [isExpanded, hasLoadedAssetTypes, onLoad]);

    useEffect(() => {
        if (isExpandedByDefault && !hasLoadedAssetTypes) {
            handleToggleAssets();
        }
    }, [isExpandedByDefault, hasLoadedAssetTypes, handleToggleAssets]);

    return (
        <div className="asset-list">
            <div className="asset-list__header">
                <button
                    onClick={handleToggleAssets}
                    disabled={isLoading}
                    className="dashboard-assets-toggle-button"
                >
                    <span className="dashboard-assets-toggle-text">
                        {isLoading ? 'Loading...' : title}
                    </span>
                    {!isLoading && (
                        <ChevronIcon className={`dashboard-assets-toggle-icon ${isExpanded ? 'rotated' : ''}`} />
                    )}
                </button>
            </div>

            {isExpanded && (
                <div className="asset-list__items">
                    {allAssets.length === 0 ? (
                        <div className="asset-list__empty-message">
                            (Empty)
                        </div>
                    ) : (
                        allAssets.map((asset, index) =>
                            asset.amount > 0 ? (
                                <InvestedAssetSummary
                                    key={`${asset.assetType}-${index}`}
                                    asset={asset}
                                    historicalValues={asset.dataPoints ?? []}
                                />
                            ) : (
                                <UninvestedAssetSummary
                                    key={`${asset.assetType}-${index}`}
                                    asset={asset}
                                    historicalValues={asset.dataPoints ?? []}
                                />
                            )
                        )
                    )}
                </div>
            )}
        </div>
    );
}
