import type { Asset } from './Asset';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';

interface AssetSummaryProps {
    asset: Asset;
}

export default function AssetSummary({ asset }: AssetSummaryProps) {
    const dataPoints = [5, 10, 5, 20, 8, 15, 10];
    const width = 32;
    const height = 20;
    const min = Math.min(...dataPoints);
    const max = Math.max(...dataPoints);
    const points = dataPoints.map((v, i) => {
        const x = dataPoints.length > 1 ? (i / (dataPoints.length - 1)) * width : 0;
        const y = max === min ? height / 2 : height - ((v - min) / (max - min)) * height;
        return `${x},${y}`;
    }).join(' ');
    return (
        <div className="asset-list__item">
            <div className="asset-list__content">
                <div className="asset-list__type">{asset.assetType}</div>
                <div className="asset-list__amount">
                    {formatCurrency(asset.amount)}
                </div>
            </div>
            <button className="asset-list__button" aria-label="View asset details">
                <svg
                    className="asset-list__spark"
                    width={width}
                    height={height}
                    viewBox={`0 0 ${width} ${height}`}
                    xmlns="http://www.w3.org/2000/svg"
                    preserveAspectRatio="none"
                    style={{ display: 'block', flex: 'none' }}
                >
                    <polyline
                        points={points}
                        fill="none"
                        stroke="#ffffff"
                        strokeWidth={2}
                        strokeLinecap="round"
                    />
                </svg>
            </button>
        </div>
    );
}
