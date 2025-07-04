import type { Asset } from './Asset';
import AssetSummary from './AssetSummary';
import './AssetList.css';

interface AssetListProps {
    assets: Asset[];
}

export default function AssetList({ assets }: AssetListProps) {
    if (assets.length === 0) {
        return (
            <div className="asset-list asset-list--empty">
                <div className="asset-list__empty-message">
                    You have no assets
                </div>
            </div>
        );
    }

    return (
        <div className="asset-list">
            <div className="asset-list__header">
                <h3 className="asset-list__title">Your Assets</h3>
            </div>
            
            <div className="asset-list__items">
                {assets.map((asset, index) => (
                    <AssetSummary key={`${asset.assetType}-${index}`} asset={asset} />
                ))}
            </div>
        </div>
    );
}
