import { useState } from 'react';
import './DashboardHeader.css';
import CapitalPopup from './CapitalPopup';
import type { PerformanceHistory } from '../User';

interface DashboardHeaderProps {
  bankName: string;
  claimedCapital: number;
  actualCapital: number;
  performanceHistory: PerformanceHistory | null;
}

type PopupType = 'claimed' | 'actual' | null;

function formatCurrency(amount: number) {
  return amount.toLocaleString(undefined, { style: 'currency', currency: 'GBP', maximumFractionDigits: 2 });
}

export default function DashboardHeader({ bankName, claimedCapital, actualCapital, performanceHistory }: DashboardHeaderProps) {
  // Performance history is available here for future use (charts, trends, etc.)
  // Currently not displayed but loaded for future features
  void performanceHistory; // Mark as intentionally unused
  
  const [activePopup, setActivePopup] = useState<PopupType>(null);

  const handleCapitalClick = (type: 'claimed' | 'actual') => {
    setActivePopup(type);
  };

  const closePopup = () => {
    setActivePopup(null);
  };

  const getPopupTitle = () => {
    return activePopup === 'claimed' ? 'Claimed Capital' : 'Actual Capital';
  };

  const getPopupValue = () => {
    return activePopup === 'claimed' ? claimedCapital : actualCapital;
  };
  
  return (
    <>
      <header className="dashboard-header">
        <div className="dashboard-header__bank">{bankName}</div>
        <div className="dashboard-header__capitals">
          <button 
            className="dashboard-header__capital dashboard-header__capital--clickable"
            onClick={() => handleCapitalClick('claimed')}
            aria-label="View claimed capital details"
          >
            <span className="dashboard-header__capital-label">Claimed Capital</span>
            <span className="dashboard-header__capital-value">
              {formatCurrency(claimedCapital)}
              <svg className="dashboard-header__capital-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </span>
          </button>
          <button 
            className="dashboard-header__capital dashboard-header__capital--clickable"
            onClick={() => handleCapitalClick('actual')}
            aria-label="View actual capital details"
          >
            <span className="dashboard-header__capital-label">Actual Capital</span>
            <span className="dashboard-header__capital-value">
              {formatCurrency(actualCapital)}
              <svg className="dashboard-header__capital-icon" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <path d="M9 18l6-6-6-6"/>
              </svg>
            </span>
          </button>
        </div>
      </header>

      <CapitalPopup
        isOpen={activePopup !== null}
        onClose={closePopup}
        title={getPopupTitle()}
        value={getPopupValue()}
        type={activePopup || 'claimed'}
      />
    </>
  );
}
