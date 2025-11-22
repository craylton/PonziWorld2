import React, { useRef } from 'react';
import { useDrag } from '@use-gesture/react';
import './BottomPanel.css';

interface BottomPanelProps {
  visible: boolean;
  children: React.ReactNode;
  onClose?: () => void;
}

const BottomPanel: React.FC<BottomPanelProps> = ({ visible, children, onClose }) => {
  const panelRef = useRef<HTMLDivElement>(null);

  const bind = useDrag(
    ({ movement: [, my], direction: [, dy], velocity: [, vy], last }) => {
      if (!visible || !onClose) return;

      if (last) {
        const swipeThreshold = 50;
        const velocityThreshold = 0.5;

        // Swipe down to close (positive movement/direction)
        const shouldClose = my > swipeThreshold || (dy > 0 && vy > velocityThreshold);

        if (shouldClose) {
          onClose();
        }
      }
    },
    {
      axis: 'y',
      filterTaps: true,
      threshold: 10,
    }
  );

  const panelClass = `dashboard-bottompanel${visible ? ' dashboard-bottompanel--visible' : ''}`;
  return (
    <div ref={panelRef} className={panelClass} {...bind()}>
      {children}
    </div>
  );
};

export default BottomPanel;
