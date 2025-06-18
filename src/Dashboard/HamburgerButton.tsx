import React from 'react';

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
    {shouldAllowClose ? (
      // X icon
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#333" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <line x1="18" y1="6" x2="6" y2="18" />
        <line x1="6" y1="6" x2="18" y2="18" />
      </svg>
    ) : (
      // Hamburger icon
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#333" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <line x1="4" y1="7" x2="20" y2="7" />
        <line x1="4" y1="12" x2="20" y2="12" />
        <line x1="4" y1="17" x2="20" y2="17" />
      </svg>
    )}
  </button>
);

export default HamburgerButton;
