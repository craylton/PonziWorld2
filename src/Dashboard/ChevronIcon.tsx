import React from 'react';

interface ChevronIconProps {
    width?: number;
    height?: number;
    stroke?: string;
    strokeWidth?: number;
    className?: string;
}

const ChevronIcon: React.FC<ChevronIconProps> = ({
    width = 16,
    height = 16,
    stroke = 'currentColor',
    strokeWidth = 2,
    className = 'dashboard-header__chevron-icon',
}
) => (
    <svg
        className={className}
        width={width}
        height={height}
        viewBox="0 0 24 24"
        fill="none"
        stroke={stroke}
        strokeWidth={strokeWidth}
    >
        <path d="M9 18l6-6-6-6" />
    </svg>
);

export default ChevronIcon;
