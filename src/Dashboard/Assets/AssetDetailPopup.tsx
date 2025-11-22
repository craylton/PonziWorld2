import { useEffect, useState, useCallback } from 'react';
import LineGraph from './LineGraph';
import { formatCurrency } from '../../utils/currency';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';
import { useAssetContext } from '../../contexts/useAssetContext';
import { useLoadingContext } from '../../contexts/useLoadingContext';
import TransactionPopup from './TransactionPopup';
import Popup from '../../components/Popup';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';
import type { HistoricalPerformanceEntry } from '../../models/HistoricalPerformance';

interface AssetDetailPopupProps {
  isOpen: boolean;
  onClose: () => void;
  investment: InvestmentDetailsResponse;
}

export default function AssetDetailPopup({
  isOpen,
  onClose,
  investment
}: AssetDetailPopupProps) {
  const { bankId } = useBankContext();
  const { refreshBank, cashBalance } = useAssetContext();
  const { showLoadingPopup } = useLoadingContext();
  const [transactionPopupOpen, setTransactionPopupOpen] = useState(false);
  const [transactionType, setTransactionType] = useState<'buy' | 'sell'>('buy');
  const [chartData, setChartData] = useState<HistoricalPerformanceEntry[]>([]);
  const [isLoadingChart, setIsLoadingChart] = useState(false);
  const [chartError, setChartError] = useState<string | null>(null);

  // Fetch real historical performance data
  const fetchChartData = useCallback(async () => {
    if (!bankId || !investment.targetAssetId) return;

    setIsLoadingChart(true);
    setChartError(null);

    try {
      const response = await makeAuthenticatedRequest(
        `/api/historicalPerformance/asset/${investment.targetAssetId}/${bankId}`,
        {
          method: 'GET',
        }
      );

      if (response.ok) {
        const data: HistoricalPerformanceEntry[] = await response.json();
        setChartData(data);
        setChartError(null);
      } else {
        const errorData = await response.json().catch(() => ({}));
        const errorMessage = errorData.error || 'Failed to fetch historical performance data';
        console.error('Failed to fetch historical performance data:', errorMessage);
        setChartError(errorMessage);
        setChartData([]);
      }
    } catch (error) {
      const errorMessage = 'Network error occurred while fetching chart data';
      console.error('Error fetching historical performance data:', error);
      setChartError(errorMessage);
      setChartData([]);
    } finally {
      setIsLoadingChart(false);
    }
  }, [bankId, investment.targetAssetId]);

  // Fetch data when popup opens or key identifiers change
  useEffect(() => {
    if (isOpen && bankId && investment.targetAssetId) {
      // Only fetch if we don't already have data for this asset
      setChartData([]);
      setChartError(null);
      fetchChartData();
    }
  }, [isOpen, bankId, investment.targetAssetId, fetchChartData]);

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
    // Close transaction popup and show global loading popup
    setTransactionPopupOpen(false);
    showLoadingPopup('loading', 'Processing transaction...');

    try {
      // Determine the endpoint based on transaction type
      const endpoint = transactionType === 'buy' ? '/api/buy' : '/api/sell';

      const response = await makeAuthenticatedRequest(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          sourceBankId: bankId,
          targetAssetId: investment.targetAssetId,
          amount,
        }),
      });

      if (response.ok) {
        showLoadingPopup('success', 'Transaction completed successfully');

        if (refreshBank) {
          await refreshBank();
        }
      } else {
        // Error from server
        const error = await response.json();
        console.error(`${transactionType} transaction failed:`, error);
        showLoadingPopup('error', error.error || 'Transaction failed');
      }
    } catch (error) {
      // Network or other error
      console.error(`Error during ${transactionType} transaction:`, error);
      showLoadingPopup('error', 'Network error occurred');
    }
  };

  const handleTransactionClose = () => {
    setTransactionPopupOpen(false);
  };

  const hasInvestmentOrPending = investment.investedAmount > 0 || investment.pendingAmount !== 0;

  const footer = hasInvestmentOrPending ? (
    <>
      <button
        className="popup__button popup__button--success"
        onClick={handleBuy}
      >
        Buy
      </button>
      <button
        className="popup__button popup__button--danger"
        onClick={handleSell}
      >
        Sell
      </button>
    </>
  ) : (
    <button
      className="popup__button popup__button--success"
      onClick={handleBuy}
    >
      Buy
    </button>
  );

  return (
    <>
      <Popup
        isOpen={isOpen}
        title={`${investment.targetAssetName} Details`}
        onClose={onClose}
        footer={footer}
        className="asset-detail-popup"
      >
        {investment.investedAmount > 0 && (
          <div className="popup__value">
            {formatCurrency(investment.investedAmount)}
          </div>
        )}
        <div className="popup__chart">
          {isLoadingChart ? (
            <div style={{ textAlign: 'center', padding: '40px', color: '#666' }}>
              Loading chart data...
            </div>
          ) : chartError ? (
            <div style={{ textAlign: 'center', padding: '40px', color: '#d32f2f' }}>
              {chartError}
              <br />
              <button
                onClick={fetchChartData}
                style={{
                  marginTop: '10px',
                  padding: '8px 16px',
                  backgroundColor: '#1976d2',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer'
                }}
              >
                Retry
              </button>
            </div>
          ) : chartData.length > 0 ? (
            <LineGraph
              data={chartData}
              title={investment.targetAssetName}
              formatTooltip={(value) => `${value}%`}
              formatYAxisTick={(value) => `${value}%`}
            />
          ) : (
            <div style={{ textAlign: 'center', padding: '40px', color: '#666' }}>
              No historical data available
            </div>
          )}
        </div>
      </Popup>
      <TransactionPopup
        isOpen={transactionPopupOpen}
        onClose={handleTransactionClose}
        assetName={investment.targetAssetName}
        transactionType={transactionType}
        currentHoldings={investment.investedAmount + investment.pendingAmount}
        maxBuyAmount={cashBalance}
        onConfirm={handleTransactionConfirm}
      />
    </>
  );
}
