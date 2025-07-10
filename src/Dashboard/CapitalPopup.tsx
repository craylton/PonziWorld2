import { useRef, useEffect } from 'react';
import './CapitalPopup.css';
import { formatCurrency } from '../utils/currency';
import type { PerformanceHistoryEntry } from '../models/PerformanceHistory';
import LineGraph from './Assets/LineGraph';

interface CapitalPopupProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  value: number;
  performanceHistory: PerformanceHistoryEntry[] | null;
  isHistoryLoading: boolean;
}

export default function CapitalPopup({
  isOpen,
  onClose,
  title,
  value,
  performanceHistory,
  isHistoryLoading
}: CapitalPopupProps) {
  const popupRef = useRef<HTMLDivElement>(null);
  const overlayRef = useRef<HTMLDivElement>(null);

  // Format the chart data based on the performance history
  const getChartData = () => {
    if (!performanceHistory) return [];

    return performanceHistory.map(entry => ({
      day: entry.day,
      value: entry.value
    }));
  };

  const chartData = getChartData();

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
          ) : performanceHistory ? (
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
