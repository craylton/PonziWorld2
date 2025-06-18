import React from 'react';
import CloseIcon from './CloseIcon';
import HamburgerIcon from './HamburgerIcon';

interface HamburgerButtonProps {
  shouldAllowClose: boolean;
  onClick: () => void;
  ariaLabel: string;
  className?: string;
}

const HamburgerButton: React.FC<HamburgerButtonProps> = ({ shouldAllowClose, onClick, ariaLabel, className = '' }) => (
  <button
    className={`dashboard-hamburger ${className} ${shouldAllowClose ? 'dashboard-hamburger--open' : ''}`}
    aria-label={ariaLabel}
    onClick={onClick}
    type="button"
  >
    {shouldAllowClose ? <CloseIcon /> : <HamburgerIcon />}
  </button>
);

export default HamburgerButton;
