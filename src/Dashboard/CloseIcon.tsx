import React from 'react';

interface XButtonIconProps {
  width?: number;
  height?: number;
  stroke?: string;
  strokeWidth?: number;
  className?: string;
}

const CloseIcon: React.FC<XButtonIconProps> = ({
  width = 24,
  height = 24,
  stroke = '#333',
  strokeWidth = 2,
  className = '',
}) => (
  <svg
    width={width}
    height={height}
    viewBox="0 0 24 24"
    fill="none"
    stroke={stroke}
    strokeWidth={strokeWidth}
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <line x1="18" y1="6" x2="6" y2="18" />
    <line x1="6" y1="6" x2="18" y2="18" />
  </svg>
);

export default CloseIcon;
