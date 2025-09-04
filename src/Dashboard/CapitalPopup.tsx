import { useEffect, useState } from 'react';
import { formatCurrency, formatCurrencyFromString } from '../utils/currency';
import { makeAuthenticatedRequest } from '../auth';
import type { HistoricalPerformanceEntry, OwnBankHistoricalPerformance } from '../models/HistoricalPerformance';
import LineGraph from './Assets/LineGraph';
import Popup from '../components/Popup';

interface CapitalPopupProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  value: string; // Now string for arbitrary precision
  bankId: string;
  capitalType: 'claimed' | 'actual';
}

export default function CapitalPopup({
  isOpen,
  onClose,
  title,
  value,
  bankId,
  capitalType
}: CapitalPopupProps) {
  const [historicalPerformance, setHistoricalPerformance] = useState<HistoricalPerformanceEntry[] | null>(null);
  const [isHistoryLoading, setIsHistoryLoading] = useState(false);

  // Fetch historical performance data when popup opens
  useEffect(() => {
    if (!isOpen) {
      return;
    }

    setIsHistoryLoading(true);
    setHistoricalPerformance(null);

    const fetchHistoricalData = async () => {
      try {
        const response = await makeAuthenticatedRequest(`/api/historicalPerformance/ownbank/${bankId}`);
        if (response.ok) {
          const data: OwnBankHistoricalPerformance = await response.json();
          const performanceData = capitalType === 'claimed' ? data.claimedHistory : data.actualHistory;
          setHistoricalPerformance(performanceData);
        }
      } catch (error) {
        console.error('Failed to fetch historical performance:', error);
      } finally {
        setIsHistoryLoading(false);
      }
    };

    fetchHistoricalData();
  }, [isOpen, bankId, capitalType]);

  // Format the chart data based on the performance history
  const getChartData = () => {
    if (!historicalPerformance) return [];

    return historicalPerformance.map((entry: HistoricalPerformanceEntry) => ({
      day: entry.day,
      value: entry.value
    }));
  };

  const chartData = getChartData();

  const footer = (
    <button
      className="popup__button popup__button--primary"
      onClick={onClose}
    >
      Close
    </button>
  );

  return (
    <Popup
      isOpen={isOpen}
      title={title}
      onClose={onClose}
      footer={footer}
      className="capital-popup"
    >
      <div className="popup__value">
        {formatCurrencyFromString(value)}
      </div>

      {isHistoryLoading ? (
        <div className="popup__loading">
          <p>Loading chart data...</p>
        </div>
      ) : historicalPerformance ? (
        <div className="popup__chart">
          <LineGraph
            data={chartData}
            title={title}
            formatTooltip={formatCurrency}
            formatYAxisTick={formatCurrency}
          />
        </div>
      ) : (
        <div className="popup__no-data">
          <p>No chart data available.</p>
        </div>
      )}
    </Popup>
  );
}
