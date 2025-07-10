import { useRef, useEffect } from 'react';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import '../CapitalPopup.css';
import { formatCurrency } from '../../utils/currency';

interface AssetDetailPopupProps {
  isOpen: boolean;
  onClose: () => void;
  assetType: string;
}

export default function AssetDetailPopup({
  isOpen,
  onClose,
  assetType
}: AssetDetailPopupProps) {
  const popupRef = useRef<HTMLDivElement>(null);
  const overlayRef = useRef<HTMLDivElement>(null);

  // Generate dummy detailed chart data (30 days) using the same algorithm as AssetSection
  const getDummyChartData = () => {
    const data = [];
    let currentValue = 1000; // Starting value
    
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

  // Calculate Y-axis domain with padding around the data range
  const getYAxisDomain = () => {
    const data = getDummyChartData();
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

  const chartData = getDummyChartData();

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
          <div className="capital-popup__value">
            {formatCurrency(123)}
          </div>

          <div className="capital-popup__chart">
            <ResponsiveContainer width="100%" height={300}>
              <LineChart
                data={chartData}
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
