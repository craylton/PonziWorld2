import React, { useRef } from 'react';
import { useDrag } from '@use-gesture/react';

interface SidePanelProps {
  side: 'left' | 'right';
  visible: boolean;
  children: React.ReactNode;
  onClose?: () => void;
}

const SidePanel: React.FC<SidePanelProps> = ({ side, visible, children, onClose }) => {
  const panelRef = useRef<HTMLDivElement>(null);

  const bind = useDrag(
    ({ movement: [mx], direction: [dx], velocity: [vx], last }) => {
      if (!visible || !onClose) return;

      // Only handle horizontal swipes
      if (last) {
        const swipeThreshold = 50;
        const velocityThreshold = 0.5;

        // For left panel: swipe left to close (negative movement/direction)
        // For right panel: swipe right to close (positive movement/direction)
        const shouldClose = side === 'left'
          ? (mx < -swipeThreshold || (dx < 0 && vx > velocityThreshold))
          : (mx > swipeThreshold || (dx > 0 && vx > velocityThreshold));

        if (shouldClose) {
          onClose();
        }
      }
    },
    {
      axis: 'x', // Only respond to horizontal gestures
      filterTaps: true, // Ignore taps
      threshold: 10, // Minimum movement to start gesture
    }
  );

  const panelClass = `dashboard-sidepanel dashboard-sidepanel--${side}${visible ? ' dashboard-sidepanel--visible' : ''}`;
  return (
    <aside ref={panelRef} className={panelClass} {...bind()}>
      {children}
    </aside>
  );
};

export default SidePanel;
