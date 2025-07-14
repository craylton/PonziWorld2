import { useEffect, useState } from 'react';
import '../CapitalPopup.css';
import LineGraph from './LineGraph';
import { formatCurrency } from '../../utils/currency';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';
import { useAssetContext } from '../../contexts/useAssetContext';
import TransactionPopup from './TransactionPopup';
import type { Asset } from './Asset';

interface AssetDetailPopupProps {
  isOpen: boolean;
  onClose: () => void;
  asset: Asset;
  onTransactionStart: () => void;
  onTransactionComplete: (success: boolean, message: string) => void;
}

export default function AssetDetailPopup({
  isOpen,
  onClose,
  asset,
  onTransactionStart,
  onTransactionComplete
}: AssetDetailPopupProps) {
  const { bankId } = useBankContext();
  const { refreshAssets } = useAssetContext();
  const [transactionPopupOpen, setTransactionPopupOpen] = useState(false);
  const [transactionType, setTransactionType] = useState<'buy' | 'sell'>('buy');

  // Generate dummy detailed chart data (30 days) using the same algorithm as AssetSection
  const getChartData = () => {
    const data = [];
    let currentValue = 100;
    
    for (let i = 0; i < 30; i++) {
      data.push({
        day: i + 1,
        value: Math.round(currentValue)
      });
      
      // Use the same algorithm: multiply by random factor between 0.9 and 1.2
      const factor = 0.9 + Math.random() * 0.3; // 0.9 to 1.2
      currentValue = currentValue * factor;
    }
    
    return data;
  };

  const chartData = getChartData();

  // Functions for Buy/Sell actions
  const handleBuy = () => {
    setTransactionType('buy');
    setTransactionPopupOpen(true);
  };

  const handleSell = () => {
    setTransactionType('sell');
    setTransactionPopupOpen(true);
  };

  const handleTransactionConfirm = async (amount: number) => {
    // Close both popups and notify parent to show loading
    setTransactionPopupOpen(false);
    // Close the asset detail popup with a slight delay to ensure loading popup appears
    setTimeout(() => {
      onClose();
    }, 10);
    onTransactionStart();

    try {
      // Determine the endpoint based on transaction type
      const endpoint = transactionType === 'buy' ? '/api/buy' : '/api/sell';

      const response = await makeAuthenticatedRequest(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          buyerBankId: bankId,
          assetId: asset.assetTypeId,
          amount,
        }),
      });

      if (response.ok) {
        onTransactionComplete(true, 'Transaction completed successfully');
        // Refresh assets in the background
        refreshAssets();
      } else {
        // Error from server
        const error = await response.json();
        console.error(`${transactionType} transaction failed:`, error);
        onTransactionComplete(false, error.error || 'Transaction failed');
      }
    } catch (error) {
      // Network or other error
      console.error(`Error during ${transactionType} transaction:`, error);
      onTransactionComplete(false, 'Network error occurred');
    }
  };

  const handleTransactionClose = () => {
    setTransactionPopupOpen(false);
  };

  // Prevent background scrolling when open
  useEffect(() => {
    if (isOpen) document.body.style.overflow = 'hidden';
    return () => { document.body.style.overflow = 'unset'; };
  }, [isOpen]);

  if (!isOpen) return null;
  
  const hasInvestmentOrPending = asset.amount > 0 || asset.pendingAmount !== 0;

  return (
    <div
      className="capital-popup-overlay"
      onClick={e => e.target === e.currentTarget && onClose()}
      role="dialog"
      aria-modal="true"
      aria-labelledby="popup-title"
    >
      <div className="capital-popup">
        <div className="capital-popup__header">
          <h2 id="popup-title" className="capital-popup__title">{asset.assetType} Details</h2>
          <button
            className="capital-popup__close-button"
            onClick={onClose}
            aria-label="Close popup"
          >
            ×
          </button>
        </div>
        <div className="capital-popup__content">
          {asset.amount > 0 && (
            <div className="capital-popup__value">
              {formatCurrency(asset.amount)}
            </div>
          )}
          <div className="capital-popup__chart">
            <LineGraph
              data={chartData}
              title={asset.assetType}
              formatTooltip={(value) => `${value}%`}
              formatYAxisTick={(value) => `${value}%`}
            />
          </div>
        </div>
        <div className="capital-popup__footer">
          {hasInvestmentOrPending ? (
            <>
              <button
                className="capital-popup__buy-button"
                onClick={handleBuy}
              >
                Buy
              </button>
              <button
                className="capital-popup__sell-button"
                onClick={handleSell}
              >
                Sell
              </button>
            </>
          ) : (
            <button
              className="capital-popup__buy-button"
              onClick={handleBuy}
            >
              Buy
            </button>
          )}
        </div>
      </div>
      <TransactionPopup
        isOpen={transactionPopupOpen}
        onClose={handleTransactionClose}
        assetType={asset.assetType}
        transactionType={transactionType}
        currentHoldings={asset.amount + asset.pendingAmount}
        onConfirm={handleTransactionConfirm}
      />
    </div>
  );
}
