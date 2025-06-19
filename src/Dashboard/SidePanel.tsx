import React from 'react';

interface SidePanelProps {
  side: 'left' | 'right';
  visible: boolean;
  children: React.ReactNode;
}

const SidePanel: React.FC<SidePanelProps> = ({ side, visible, children }) => {
  const panelClass = `dashboard-sidepanel dashboard-sidepanel--${side}${visible ? ' dashboard-sidepanel--visible' : ''}`;
  return (
    <aside className={panelClass}>
      {children}
    </aside>
  );
};

export default SidePanel;
