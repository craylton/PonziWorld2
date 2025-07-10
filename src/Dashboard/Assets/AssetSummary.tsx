import { useState } from 'react';
import type { Asset } from './Asset';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';
import AssetSummaryChart from './AssetSummaryChart';
import AssetDetailPopup from './AssetDetailPopup';

interface AssetSummaryProps {
    asset: Asset;
    historicalValues: number[];
}

export default function AssetSummary({ asset, historicalValues }: AssetSummaryProps) {
    const [isPopupOpen, setIsPopupOpen] = useState(false);

    // Calculate 1-day percentage change (between last two data points)
    const oneDayChange = historicalValues.length >= 2 
        ? ((historicalValues[historicalValues.length - 1] - historicalValues[historicalValues.length - 2]) / historicalValues[historicalValues.length - 2]) * 100
        : 0;

    // Calculate 7-day percentage change (between current and 7 days ago)
    const sevenDayChange = historicalValues.length >= 8
        ? ((historicalValues[historicalValues.length - 1] - historicalValues[historicalValues.length - 8]) / historicalValues[historicalValues.length - 8]) * 100
        : 0;

    // Helper function to get percentage class
    const getPercentageClass = (change: number) => {
        if (change > 0) return 'asset-list__performance-percentage--positive';
        if (change < 0) return 'asset-list__performance-percentage--negative';
        return 'asset-list__performance-percentage--neutral';
    };

    const handleChartClick = () => {
        setIsPopupOpen(true);
    };

    const handleClosePopup = () => {
        setIsPopupOpen(false);
    };

    return (
        <>
            <div className="asset-list__item">
                <div className="asset-list__content">
                    <div className="asset-list__type">{asset.assetType}</div>
                    <div className="asset-list__amount">
                        {formatCurrency(asset.amount)}
                    </div>
                </div>
                <div className="asset-list__chart-section">
                    <button 
                        className="asset-list__button" 
                        aria-label="View asset details"
                        onClick={handleChartClick}
                    >
                        <AssetSummaryChart historicalValues={historicalValues} />
                    </button>
                    <div className="asset-list__performance">
                        <div className="asset-list__performance-item">
                            1d: <span className={`asset-list__performance-percentage ${getPercentageClass(oneDayChange)}`}>
                                {oneDayChange >= 0 ? '+' : ''}{oneDayChange.toFixed(1)}%
                            </span>
                        </div>
                        <div className="asset-list__performance-item">
                            7d: <span className={`asset-list__performance-percentage ${getPercentageClass(sevenDayChange)}`}>
                                {sevenDayChange >= 0 ? '+' : ''}{sevenDayChange.toFixed(1)}%
                            </span>
                        </div>
                    </div>
                </div>
            </div>
            
            <AssetDetailPopup
                isOpen={isPopupOpen}
                onClose={handleClosePopup}
                assetType={asset.assetType}
            />
        </>
    );
}
