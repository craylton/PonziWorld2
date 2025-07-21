import { useEffect, useState } from 'react';
import './CapitalPopup.css';
import { formatCurrency } from '../utils/currency';
import { makeAuthenticatedRequest } from '../auth';
import type { HistoricalPerformanceEntry, OwnBankHistoricalPerformance } from '../models/HistoricalPerformance';
import LineGraph from './Assets/LineGraph';

interface CapitalPopupProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  value: number;
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
          <h2 id="popup-title" className="capital-popup__title">{title}</h2>
          <button
            className="capital-popup__close-button"
            onClick={onClose}
            aria-label="Close popup"
          >
            ×
          </button>
        </div>
        <div className="capital-popup__content">
          <div className="capital-popup__value">
            {formatCurrency(value)}
          </div>

          {isHistoryLoading ? (
            <div className="capital-popup__loading">
              <p>Loading chart data...</p>
            </div>
          ) : historicalPerformance ? (
            <div className="capital-popup__chart">
              <LineGraph
                data={chartData}
                title={title}
                formatTooltip={formatCurrency}
                formatYAxisTick={formatCurrency}
              />
            </div>
          ) : (
            <div className="capital-popup__no-data">
              <p>No chart data available.</p>
            </div>
          )}
        </div>
        <div className="capital-popup__footer">
          <button
            className="capital-popup__close-footer-button"
            onClick={onClose}
          >
            Close
          </button>
        </div>
      </div>
    </div>
  );
}
