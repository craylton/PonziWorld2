import { useState } from 'react';
import type { Asset } from './Asset';
import './AssetList.css';
import AssetSummaryChart from './AssetSummaryChart';
import AssetDetailPopup from './AssetDetailPopup';

interface AssetSummaryBaseProps {
  asset: Asset;
  historicalValues: number[];
  isInvested?: boolean;
}

export default function AssetSummaryBase({ asset, historicalValues, isInvested = false }: AssetSummaryBaseProps) {
  const [isPopupOpen, setIsPopupOpen] = useState(false);

  const oneDayChange = historicalValues.length >= 2
    ? ((historicalValues[historicalValues.length - 1] - historicalValues[historicalValues.length - 2])
      / historicalValues[historicalValues.length - 2]) * 100
    : 0;

  const sevenDayChange = historicalValues.length >= 8
    ? ((historicalValues[historicalValues.length - 1] - historicalValues[historicalValues.length - 8])
      / historicalValues[historicalValues.length - 8]) * 100
    : 0;

  const getPercentageClass = (change: number) => {
    if (change > 0) return 'asset-list__performance-percentage--positive';
    if (change < 0) return 'asset-list__performance-percentage--negative';
    return 'asset-list__performance-percentage--neutral';
  };

  const handleChartClick = () => setIsPopupOpen(true);
  const handleClosePopup = () => setIsPopupOpen(false);

  return (
    <>
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
      <AssetDetailPopup
        isOpen={isPopupOpen}
        onClose={handleClosePopup}
        assetType={asset.assetType}
        isInvested={isInvested}
        investedAmount={isInvested ? asset.amount : undefined}
      />
    </>
  );
}
