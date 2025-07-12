import { useEffect, useState } from 'react';
import '../CapitalPopup.css';
import LineGraph from './LineGraph';
import { formatCurrency } from '../../utils/currency';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';
import TransactionPopup from './TransactionPopup';

interface AssetDetailPopupProps {
  isOpen: boolean;
  onClose: () => void;
  assetType: string;
  assetTypeId: string;
  investedAmount: number;
}

export default function AssetDetailPopup({
  isOpen,
  onClose,
  assetType,
  assetTypeId,
  investedAmount
}: AssetDetailPopupProps) {
  const { bankId } = useBankContext();
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
    try {
      // Determine the endpoint and amount based on transaction type
      const endpoint = transactionType === 'buy' ? '/api/buy' : '/api/sell';
      const finalAmount = transactionType === 'buy' ? amount : -amount; // Negative for sell

      const response = await makeAuthenticatedRequest(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          buyerBankId: bankId,
          assetId: assetTypeId,
          amount: finalAmount,
        }),
      });

      if (response.ok) {
        const result = await response.json();
        console.log(`${transactionType} transaction successful:`, result);
        // You could add a success notification here
        // For now, just close the popup
        setTransactionPopupOpen(false);
      } else {
        const error = await response.json();
        console.error(`${transactionType} transaction failed:`, error);
        // You could add an error notification here
      }
    } catch (error) {
      console.error(`Error during ${transactionType} transaction:`, error);
      // You could add an error notification here
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
          {investedAmount > 0 && (
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
          {investedAmount > 0 ? (
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
