import { useState, useEffect } from 'react';
import type { Asset } from './Asset';
import type { AssetDetailsResponse } from '../../models/AssetDetails';
import './AssetList.css';
import AssetSummaryChart from './AssetSummaryChart';
import AssetDetailPopup from './AssetDetailPopup';
import LoadingPopup from './LoadingPopup';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';

interface AssetSummaryBaseProps {
  asset: Asset;
  onAssetDetailsLoaded?: (investedAmount: number, pendingAmount: number) => void;
}

export default function AssetSummaryBase({ asset, onAssetDetailsLoaded }: AssetSummaryBaseProps) {
  const { bankId } = useBankContext();
  const [isPopupOpen, setIsPopupOpen] = useState(false);
  const [loadingPopupOpen, setLoadingPopupOpen] = useState(false);
  const [loadingStatus, setLoadingStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [loadingMessage, setLoadingMessage] = useState<string>('');
  const [assetDetails, setAssetDetails] = useState<AssetDetailsResponse | null>(null);
  const [isLoadingData, setIsLoadingData] = useState(false);

  // Fetch asset details including historical data
  useEffect(() => {
    const fetchAssetDetails = async () => {
      setIsLoadingData(true);
      try {
        const response = await makeAuthenticatedRequest(`/api/asset/${asset.assetTypeId}/${bankId}`);
        if (response.ok) {
          const data: AssetDetailsResponse = await response.json();
          setAssetDetails(data);
          
          // Notify parent component about the loaded asset details
          if (onAssetDetailsLoaded) {
            onAssetDetailsLoaded(data.investedAmount, data.pendingAmount);
          }
        } else {
          console.error('Failed to fetch asset details for asset:', asset.assetType);
        }
      } catch (error) {
        console.error('Error fetching asset details:', error);
      } finally {
        setIsLoadingData(false);
      }
    };

    fetchAssetDetails();
  }, [asset.assetTypeId, bankId, asset.assetType, onAssetDetailsLoaded]);

  // Use historical data values for calculations, or empty array if loading
  const historicalValues = assetDetails?.historicalData?.map(entry => entry.value) || [];
  
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
          disabled={isLoadingData}
        >
          {isLoadingData ? (
            <div className="asset-list__loading">Loading...</div>
          ) : (
            <AssetSummaryChart historicalValues={historicalValues} />
          )}
        </button>
        <div className="asset-list__performance">
          <div className="asset-list__performance-item">
            1d: <span className={`asset-list__performance-percentage ${getPercentageClass(oneDayChange)}`}>
              {isLoadingData ? '--' : `${oneDayChange >= 0 ? '+' : ''}${oneDayChange.toFixed(1)}%`}
            </span>
          </div>
          <div className="asset-list__performance-item">
            7d: <span className={`asset-list__performance-percentage ${getPercentageClass(sevenDayChange)}`}>
              {isLoadingData ? '--' : `${sevenDayChange >= 0 ? '+' : ''}${sevenDayChange.toFixed(1)}%`}
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
