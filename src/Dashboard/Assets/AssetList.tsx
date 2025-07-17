import './AssetList.css';
import ChevronIcon from '../ChevronIcon';
import { useState, useEffect, useCallback } from 'react';
import { useAssetContext } from '../../contexts/useAssetContext';
import type { AvailableAsset } from '../../models/AvailableAsset';
import AssetSummary from './AssetSummary';

interface AssetListProps {
    title: string;
    onLoad: () => Promise<AvailableAsset[]>;
    isExpandedByDefault: boolean;
}

export default function AssetList({ title, onLoad, isExpandedByDefault }: AssetListProps) {
    const [allAssets, setAllAssets] = useState<AvailableAsset[]>([]);
    const [isExpanded, setIsExpanded] = useState(false);
    const [isLoading, setIsLoading] = useState(false);
    const [hasLoadedAssetTypes, setHasLoadedAssetTypes] = useState(false);
    const { registerRefreshCallback, unregisterRefreshCallback } = useAssetContext();

    const loadAssets = useCallback(async () => {
        setIsLoading(true);
        try {
            const assets = await onLoad();
            setAllAssets(assets);
            setHasLoadedAssetTypes(true);
        } catch (error) {
            console.error('Error loading asset types:', error);
        } finally {
            setIsLoading(false);
        }
    }, [onLoad]);

    const handleRefresh = useCallback(async () => {
        if (hasLoadedAssetTypes && isExpanded) {
            await loadAssets();
        }
    }, [hasLoadedAssetTypes, isExpanded, loadAssets]);

    // Register for refresh callbacks
    useEffect(() => {
        registerRefreshCallback(handleRefresh);
        return () => {
            unregisterRefreshCallback(handleRefresh);
        };
    }, [registerRefreshCallback, unregisterRefreshCallback, handleRefresh]);

    const handleToggleAssets = useCallback(async () => {
        if (isExpanded) {
            setIsExpanded(false);
            return;
        }

        if (hasLoadedAssetTypes) {
            setIsExpanded(true);
            return;
        }

        await loadAssets();
        setIsExpanded(true);
    }, [isExpanded, hasLoadedAssetTypes, loadAssets]);

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
                        allAssets.map((asset, index) => {
                            return (
                                <AssetSummary
                                    key={`${asset.assetType}-${index}`}
                                    availableAsset={asset}
                                />
                            );
                        })
                    )}
                </div>
            )}
        </div>
    );
}
