import { useState } from 'react';
import './DashboardHeader.css';
import CapitalPopup from './CapitalPopup';
import ChevronIcon from './ChevronIcon';
import { formatCurrencyFromString } from '../utils/currency';

interface DashboardHeaderProps {
  currentDay: number;
  bankName: string;
  claimedCapital: string; // Now string for arbitrary precision
  actualCapital: string;  // Now string for arbitrary precision
  bankId: string;
}

type PopupType = 'claimed' | 'actual' | null;

export default function DashboardHeader({
  currentDay,
  bankName,
  claimedCapital,
  actualCapital,
  bankId
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

  return (
    <>
      <header className="dashboard-header">
        <div className="dashboard-header__day">Day {currentDay}</div>
        <div className="dashboard-header__bank">{bankName}</div>
        <div className="dashboard-header__capitals">
          <button
            className="dashboard-header__capital dashboard-header__capital--clickable"
            onClick={() => handleCapitalClick('claimed')}
            aria-label="View claimed capital details"
          >
            <span className="dashboard-header__capital-label">Claimed Capital</span>
            <span className="dashboard-header__capital-value">
              {formatCurrencyFromString(claimedCapital)}
              <ChevronIcon />
            </span>
          </button>
          <button
            className="dashboard-header__capital dashboard-header__capital--clickable"
            onClick={() => handleCapitalClick('actual')}
            aria-label="View actual capital details"
          >
            <span className="dashboard-header__capital-label">Actual Capital</span>
            <span className="dashboard-header__capital-value">
              {formatCurrencyFromString(actualCapital)}
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
        bankId={bankId}
        capitalType={activePopup || 'claimed'}
      />
    </>
  );
}
