import type { Asset } from './Asset';
import AssetSummary from './AssetSummary';
import './AssetList.css';
import ChevronIcon from '../ChevronIcon';
import { useState, useEffect, useCallback } from 'react';

interface AssetListProps {
    title: string;
    onLoad: () => Promise<Asset[]>;
    isExpandedByDefault: boolean;
}

export default function AssetList({ title, onLoad, isExpandedByDefault }: AssetListProps) {
    const [allAssets, setAllAssets] = useState<Asset[]>([]);
    const [showAssets, setShowAssets] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [hasLoadedAssetTypes, setHasLoadedAssetTypes] = useState(false);

    const handleToggleAssets = useCallback(async () => {
        if (showAssets) {
            setShowAssets(false);
            return;
        }

        if (hasLoadedAssetTypes) {
            setShowAssets(true);
            return;
        }

        setIsLoading(true);
        try {
            setAllAssets(await onLoad());
            setHasLoadedAssetTypes(true);
            setShowAssets(true);
        } catch (error) {
            console.error('Error loading asset types:', error);
        } finally {
            setIsLoading(false);
        }
    }, [showAssets, hasLoadedAssetTypes, onLoad]);

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
                        <ChevronIcon className={`dashboard-assets-toggle-icon ${showAssets ? 'rotated' : ''}`} />
                    )}
                </button>
            </div>

            {showAssets && (
                <div className="asset-list__items">
                    {allAssets.length === 0 ? (
                        <div className="asset-list__empty-message">
                            (Empty)
                        </div>
                    ) : (
                        allAssets.map((asset, index) => (
                            <AssetSummary key={`${asset.assetType}-${index}`} asset={asset} />
                        ))
                    )}
                </div>
            )}
        </div>
    );
}
