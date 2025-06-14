import React from 'react';

interface SidePanelProps {
  side: 'left' | 'right';
  visible: boolean;
  children: React.ReactNode;
}

const SidePanel: React.FC<SidePanelProps> = ({ side, visible, children }) => {
  const panelClass = `dashboard-panel dashboard-panel--${side}${visible ? ' dashboard-panel--visible' : ''}`;
  return (
    <aside className={panelClass}>
      {children}
    </aside>
  );
};

export default SidePanel;
