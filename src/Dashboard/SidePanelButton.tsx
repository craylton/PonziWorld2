import React from 'react';
import CloseIcon from './CloseIcon';
import HamburgerIcon from './HamburgerIcon';
import CogIcon from './CogIcon';

interface SidePanelButtonProps {
    iconType: 'cog' | 'hamburger';
    shouldAllowClose: boolean;
    onClick: () => void;
    ariaLabel: string;
    className?: string;
}

const SidePanelButton: React.FC<SidePanelButtonProps> = ({
    iconType,
    shouldAllowClose,
    onClick,
    ariaLabel,
    className = '' }) => (
    <button
        className={`dashboard-sidepanel-button ${className}`}
        aria-label={ariaLabel}
        onClick={onClick}
        type="button"
    >
        {shouldAllowClose ? <CloseIcon /> : iconType === 'cog' ? <CogIcon /> : <HamburgerIcon />}
    </button>
);

export default SidePanelButton;
