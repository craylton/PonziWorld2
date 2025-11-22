import './AssetList.css';
import ChevronIcon from '../ChevronIcon';
import { useState, useCallback } from 'react';
import type { AssetWithDetails } from '../../models/AssetWithDetails';
import AssetSummary from './AssetSummary';

interface AssetListProps {
    title: string;
    isExpandedByDefault: boolean;
    assets: AssetWithDetails[];
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
                        assets.map((assetWithDetails, index) => {
                            return (
                                <AssetSummary
                                    key={`${assetWithDetails.asset.assetName}-${index}`}
                                    availableAsset={assetWithDetails.asset}
                                    investmentDetails={assetWithDetails.details}
                                />
                            );
                        })
                    )}
                </div>
            )}
        </div>
    );
}
