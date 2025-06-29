import { useRef, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import './CapitalPopup.css';
import { formatCurrency } from '../utils/currency';
import type { PerformanceHistoryEntry } from '../models/PerformanceHistory';

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

  // Calculate Y-axis domain with padding around the data range
  const getYAxisDomain = () => {
    const data = getChartData();
    if (data.length === 0) return ['auto', 'auto'];

    const values = data.map(d => d.value);
    const min = Math.min(...values);
    const max = Math.max(...values);

    // If all values are the same, add some padding
    if (min === max) {
      const padding = Math.max(min * 0.1, 50); // 10% padding or minimum 50
      return [min - padding, max + padding];
    }

    // Add 5% padding on each side of the range
    const range = max - min;
    const padding = range * 0.05;

    return [min - padding, max + padding];
  };

  // Custom tooltip formatter for currency
  const CustomTooltip = ({ active, payload, label }: {
    active?: boolean;
    payload?: Array<{ value: number }>;
    label?: string | number
  }) => {
    if (active && payload && payload.length) {
      return (
        <div className="capital-popup__chart-tooltip">
          <p className="label">{`Day ${label}`}</p>
          <p className="value">{formatCurrency(payload[0].value)}</p>
        </div>
      );
    }
    return null;
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
              <ResponsiveContainer width="100%" height={300}>
                <LineChart
                  data={getChartData()}
                  margin={{
                    top: 5,
                    right: 30,
                    left: 20,
                    bottom: 5,
                  }}
                >
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis
                    dataKey="day"
                    type="number"
                    domain={['dataMin', 'dataMax']}
                    tickFormatter={(value) => value.toString()}
                  />
                  <YAxis
                    domain={getYAxisDomain()}
                    tickFormatter={(value) => formatCurrency(value)}
                  />
                  <Tooltip content={<CustomTooltip />} />
                  <Line
                    type="monotone"
                    dataKey="value"
                    stroke="#2563eb"
                    strokeWidth={2}
                    dot={{ fill: '#2563eb', strokeWidth: 2, r: 3 }}
                    activeDot={{ r: 5 }}
                  />
                </LineChart>
              </ResponsiveContainer>
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
