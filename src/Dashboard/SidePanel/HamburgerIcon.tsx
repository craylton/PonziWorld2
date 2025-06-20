import React from 'react';

interface HamburgerIconProps {
    width?: number;
    height?: number;
    stroke?: string;
    strokeWidth?: number;
    className?: string;
}

const HamburgerIcon: React.FC<HamburgerIconProps> = ({
    width = 24,
    height = 24,
    stroke = '#333',
    strokeWidth = 2,
    className = '',
}
) => (
    <svg
        width={width}
        height={height}
        viewBox="0 0 24 24"
        fill="none"
        stroke={stroke}
        strokeWidth={strokeWidth}
        strokeLinecap="round"
        strokeLinejoin="round"
        className={className}>
        <line x1="4" y1="7" x2="20" y2="7" />
        <line x1="4" y1="12" x2="20" y2="12" />
        <line x1="4" y1="17" x2="20" y2="17" />
    </svg>
);

export default HamburgerIcon;
