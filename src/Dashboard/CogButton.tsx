import React from 'react';
import CloseIcon from './CloseIcon';
import CogIcon from './CogIcon';

interface CogButtonProps {
  shouldAllowClose: boolean;
  onClick: () => void;
  ariaLabel: string;
  className?: string;
}

const CogButton: React.FC<CogButtonProps> = ({ shouldAllowClose, onClick, ariaLabel, className = '' }) => (
  <button
    className={`dashboard-hamburger ${className} ${shouldAllowClose ? 'dashboard-hamburger--open' : ''}`}
    aria-label={ariaLabel}
    onClick={onClick}
    type="button"
  >
    {shouldAllowClose ? <CloseIcon /> : <CogIcon />}
  </button>
);

export default CogButton;
