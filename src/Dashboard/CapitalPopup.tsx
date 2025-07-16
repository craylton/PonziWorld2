import { useEffect } from 'react';
import './CapitalPopup.css';
import { formatCurrency } from '../utils/currency';
import type { HistoricalPerformanceEntry } from '../models/HistoricalPerformance';
import LineGraph from './Assets/LineGraph';

interface CapitalPopupProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  value: number;
  historicalPerformance: HistoricalPerformanceEntry[] | null;
  isHistoryLoading: boolean;
}

export default function CapitalPopup({
  isOpen,
  onClose,
  title,
  value,
  historicalPerformance: historicalPerformance,
  isHistoryLoading
}: CapitalPopupProps) {
  // Format the chart data based on the performance history
  const getChartData = () => {
    if (!historicalPerformance) return [];

    return historicalPerformance.map(entry => ({
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
