import { useRef, useEffect, useState } from 'react';
import '../CapitalPopup.css';
import LineGraph from './LineGraph';
import { formatCurrency } from '../../utils/currency';
import TransactionPopup from './TransactionPopup';

interface AssetDetailPopupProps {
  isOpen: boolean;
  onClose: () => void;
  assetType: string;
  isInvested?: boolean;
  investedAmount?: number;
}

export default function AssetDetailPopup({
  isOpen,
  onClose,
  assetType,
  isInvested = false,
  investedAmount
}: AssetDetailPopupProps) {
  const popupRef = useRef<HTMLDivElement>(null);
  const overlayRef = useRef<HTMLDivElement>(null);
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

  const handleTransactionConfirm = (amount: number) => {
    // TODO: Implement actual buy/sell functionality
    console.log(`${transactionType} ${amount} of ${assetType}`);
  };

  const handleTransactionClose = () => {
    setTransactionPopupOpen(false);
  };

  // Close popup when clicking outside
  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (overlayRef.current && event.target === overlayRef.current) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener('mousedown', handleClickOutside);
      document.body.style.overflow = 'hidden'; // Prevent background scrolling
    }

    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
      document.body.style.overflow = 'unset';
    };
  }, [isOpen, onClose]);

  if (!isOpen) return null;

  return (
    <div
      className="capital-popup-overlay"
      ref={overlayRef}
      role="dialog"
      aria-modal="true"
      aria-labelledby="popup-title"
    >
      <div className="capital-popup" ref={popupRef}>
        <div className="capital-popup__header">
          <h2 id="popup-title" className="capital-popup__title">{assetType} Details</h2>
          <button
            className="capital-popup__close-button"
            onClick={onClose}
            aria-label="Close popup"
          >
            ×
          </button>
        </div>
        <div className="capital-popup__content">
          {isInvested && investedAmount && (
            <div className="capital-popup__value">
              {formatCurrency(investedAmount)}
            </div>
          )}
          <div className="capital-popup__chart">
            <LineGraph
              data={chartData}
              title={assetType}
              formatTooltip={(value) => `${value}%`}
              formatYAxisTick={(value) => `${value}%`}
            />
          </div>
        </div>
        <div className="capital-popup__footer">
          {isInvested ? (
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
        assetType={assetType}
        transactionType={transactionType}
        currentHoldings={investedAmount || 0}
        onConfirm={handleTransactionConfirm}
      />
    </div>
  );
}
