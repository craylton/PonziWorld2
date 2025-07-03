import { useState } from 'react';
import './DashboardHeader.css';
import CapitalPopup from './CapitalPopup';
import ChevronIcon from './ChevronIcon';
import { formatCurrency } from '../utils/currency';
import type { PerformanceHistory } from '../models/PerformanceHistory';

interface DashboardHeaderProps {
  currentDay: number;
  bankName: string;
  claimedCapital: number;
  actualCapital: number;
  performanceHistory: PerformanceHistory | null;
  isHistoryLoading: boolean;
}

type PopupType = 'claimed' | 'actual' | null;

export default function DashboardHeader({
  currentDay,
  bankName,
  claimedCapital,
  actualCapital,
  performanceHistory,
  isHistoryLoading
}: DashboardHeaderProps) {
  const [activePopup, setActivePopup] = useState<PopupType>(null);

  const handleCapitalClick = (type: 'claimed' | 'actual') => {
    setActivePopup(type);
  };

  const handleClosePopup = () => {
    setActivePopup(null);
  };

  const getPopupTitle = () => {
    return activePopup === 'claimed' ? 'Claimed Capital' : 'Actual Capital';
  };

  const getPopupValue = () => {
    return activePopup === 'claimed' ? claimedCapital : actualCapital;
  };

  const getPopupPerformanceHistory = () => {
    if (!performanceHistory) return [];

    return activePopup === 'claimed'
      ? performanceHistory.claimedHistory
      : performanceHistory.actualHistory;
  };

  return (
    <>
      <header className="dashboard-header">
        <div className="dashboard-header__day">Day {currentDay}</div>
        <div className="dashboard-header__bank">{bankName}</div>
        <div className="dashboard-header__capitals">
          <button
            className={`dashboard-header__capital dashboard-header__capital--clickable ${isHistoryLoading ? 'dashboard-header__capital--loading' : ''}`}
            onClick={() => handleCapitalClick('claimed')}
            aria-label="View claimed capital details"
            disabled={isHistoryLoading}
          >
            <span className="dashboard-header__capital-label">Claimed Capital</span>
            <span className="dashboard-header__capital-value">
              {formatCurrency(claimedCapital)}
              <ChevronIcon />
            </span>
          </button>
          <button
            className={`dashboard-header__capital dashboard-header__capital--clickable ${isHistoryLoading ? 'dashboard-header__capital--loading' : ''}`}
            onClick={() => handleCapitalClick('actual')}
            aria-label="View actual capital details"
            disabled={isHistoryLoading}
          >
            <span className="dashboard-header__capital-label">Actual Capital</span>
            <span className="dashboard-header__capital-value">
              {formatCurrency(actualCapital)}
              <ChevronIcon />
            </span>
          </button>
        </div>
      </header>

      <CapitalPopup
        isOpen={activePopup !== null}
        onClose={handleClosePopup}
        title={getPopupTitle()}
        value={getPopupValue()}
        performanceHistory={getPopupPerformanceHistory()}
        isHistoryLoading={isHistoryLoading}
      />
    </>
  );
}
