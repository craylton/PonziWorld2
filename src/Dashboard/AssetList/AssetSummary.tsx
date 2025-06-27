import type { Asset } from './Asset';
import './AssetList.css';

interface AssetSummaryProps {
    asset: Asset;
}

export default function AssetSummary({ asset }: AssetSummaryProps) {
    const formatCurrency = (amount: number) => {
        return new Intl.NumberFormat('en-GB', {
            style: 'currency',
            currency: 'GBP',
            minimumFractionDigits: 0,
            maximumFractionDigits: 0,
        }).format(amount);
    };

    return (
        <div className="asset-list__item">
            <div className="asset-list__content">
                <div className="asset-list__type">{asset.assetType}</div>
                <div className="asset-list__amount">
                    {formatCurrency(asset.amount)}
                </div>
            </div>
            <button className="asset-list__button" aria-label="View asset details">
                📈
            </button>
        </div>
    );
}
