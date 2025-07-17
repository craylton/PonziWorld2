import './AssetList.css';
import ChevronIcon from '../ChevronIcon';
import { useState, useCallback } from 'react';
import type { AvailableAsset } from '../../models/AvailableAsset';
import AssetSummary from './AssetSummary';

interface AssetListProps {
    title: string;
    isExpandedByDefault: boolean;
    assets: AvailableAsset[];
}

export default function AssetList({ title, isExpandedByDefault, assets }: AssetListProps) {
    const [isExpanded, setIsExpanded] = useState(isExpandedByDefault);

    const handleToggleAssets = useCallback(() => {
        setIsExpanded(prev => !prev);
    }, []);

    return (
        <div className="asset-list">
            <div className="asset-list__header">
                <button
                    onClick={handleToggleAssets}
                    className="dashboard-assets-toggle-button"
                >
                    <span className="dashboard-assets-toggle-text">
                        {title}
                    </span>
                    <ChevronIcon className={`dashboard-assets-toggle-icon ${isExpanded ? 'rotated' : ''}`} />
                </button>
            </div>

            {isExpanded && (
                <div className="asset-list__items">
                    {assets.length === 0 ? (
                        <div className="asset-list__empty-message">
                            (Empty)
                        </div>
                    ) : (
                        assets.map((asset, index) => {
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
