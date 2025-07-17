import { useState } from 'react';
import type { AssetDetailsResponse } from '../../models/AssetDetails';
import './AssetList.css';
import AssetSummaryChart from './AssetSummaryChart';
import AssetDetailPopup from './AssetDetailPopup';
import LoadingPopup from './LoadingPopup';

interface AssetSummaryBaseProps {
  asset: AssetDetailsResponse;
  skipDataFetch?: boolean; // New prop to skip data fetching when AssetSummary already provides it
}

export default function AssetSummaryBase({ asset }: AssetSummaryBaseProps) {
  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [loadingPopupOpen, setLoadingPopupOpen] = useState(false);
  const [loadingStatus, setLoadingStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [loadingMessage, setLoadingMessage] = useState<string>('');

  // Use historical data values for calculations, or empty array if loading
  const historicalValues = asset?.historicalData?.map(entry => entry.value) || [];

  const numValues = historicalValues.length;
  const oneDayChange = numValues >= 2
    ? ((historicalValues[numValues - 1] - historicalValues[numValues - 2])
      / historicalValues[numValues - 2]) * 100
    : 0;

  const sevenDayChange = numValues >= 8
    ? ((historicalValues[numValues - 1] - historicalValues[numValues - 8])
      / historicalValues[numValues - 8]) * 100
    : 0;

  const getPercentageClass = (change: number) => {
    if (change > 0) return 'asset-list__performance-percentage--positive';
    if (change < 0) return 'asset-list__performance-percentage--negative';
    return 'asset-list__performance-percentage--neutral';
  };

  const handleChartClick = () => setIsPopupOpen(true);
  const handleClosePopup = () => setIsPopupOpen(false);

  const handleTransactionStart = () => {
    setLoadingPopupOpen(true);
    setLoadingStatus('loading');
    setLoadingMessage('');
  };

  const handleTransactionComplete = (success: boolean, message: string) => {
    setLoadingStatus(success ? 'success' : 'error');
    setLoadingMessage(message);
  };

  const handleLoadingClose = () => {
    setLoadingPopupOpen(false);
  };

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
        asset={asset}
        onTransactionStart={handleTransactionStart}
        onTransactionComplete={handleTransactionComplete}
      />
      <LoadingPopup
        isOpen={loadingPopupOpen}
        onClose={handleLoadingClose}
        status={loadingStatus}
        message={loadingMessage}
      />
    </>
  );
}
