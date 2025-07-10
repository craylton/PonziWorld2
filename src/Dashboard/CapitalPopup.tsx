import { useRef, useEffect } from 'react';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import { Line } from 'react-chartjs-2';
import './CapitalPopup.css';
import { formatCurrency } from '../utils/currency';
import type { PerformanceHistoryEntry } from '../models/PerformanceHistory';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

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

  // Prepare chart data for Chart.js
  const chartData = getChartData();
  const data = {
    labels: chartData.map(d => `Day ${d.day}`),
    datasets: [
      {
        label: title,
        data: chartData.map(d => d.value),
        borderColor: '#2563eb',
        backgroundColor: 'rgba(37, 99, 235, 0.1)',
        borderWidth: 2,
        pointBackgroundColor: '#2563eb',
        pointBorderColor: '#2563eb',
        pointRadius: 3,
        pointHoverRadius: 5,
        tension: 0.1,
      },
    ],
  };

  const options = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
      },
      tooltip: {
        callbacks: {
          label: (context: { parsed: { y: number } }) => {
            return formatCurrency(context.parsed.y);
          },
        },
      },
    },
    scales: {
      x: {
        grid: {
          color: 'rgba(0, 0, 0, 0.1)',
        },
        ticks: {
          maxTicksLimit: 7,
        },
      },
      y: {
        grid: {
          color: 'rgba(0, 0, 0, 0.1)',
        },
        ticks: {
          callback: (value: number | string) => {
            return formatCurrency(value as number);
          },
          maxTicksLimit: 6,
        },
      },
    },
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
              <div style={{ height: '300px' }}>
                <Line data={data} options={options} />
              </div>
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
