import React, { useRef, useEffect } from 'react';

interface SidePanelProps {
  side: 'left' | 'right';
  visible: boolean;
  children: React.ReactNode;
  onClose?: () => void;
}

const SidePanel: React.FC<SidePanelProps> = ({ side, visible, children, onClose }) => {
  const panelRef = useRef<HTMLDivElement>(null);
  const touchStartX = useRef<number>(0);
  const touchStartY = useRef<number>(0);
  const isDragging = useRef<boolean>(false);

  useEffect(() => {
    if (!visible || !onClose) return;

    const handleTouchStart = (e: TouchEvent) => {
      touchStartX.current = e.touches[0].clientX;
      touchStartY.current = e.touches[0].clientY;
      isDragging.current = false;
    };

    const handleTouchMove = (e: TouchEvent) => {
      if (!isDragging.current) {
        const deltaX = Math.abs(e.touches[0].clientX - touchStartX.current);
        const deltaY = Math.abs(e.touches[0].clientY - touchStartY.current);
        
        // Only start dragging if horizontal movement is more significant than vertical
        if (deltaX > deltaY && deltaX > 10) {
          isDragging.current = true;
        }
      }
    };

    const handleTouchEnd = (e: TouchEvent) => {
      if (!isDragging.current) return;

      const touchEndX = e.changedTouches[0].clientX;
      const deltaX = touchEndX - touchStartX.current;
      const swipeThreshold = 50; // Minimum distance for a swipe

      // For left panel: swipe left to close (negative delta)
      // For right panel: swipe right to close (positive delta)
      const shouldClose = side === 'left' ? deltaX < -swipeThreshold : deltaX > swipeThreshold;

      if (shouldClose) {
        onClose();
      }

      isDragging.current = false;
    };

    const panel = panelRef.current;
    if (panel) {
      panel.addEventListener('touchstart', handleTouchStart, { passive: true });
      panel.addEventListener('touchmove', handleTouchMove, { passive: true });
      panel.addEventListener('touchend', handleTouchEnd, { passive: true });
    }

    return () => {
      if (panel) {
        panel.removeEventListener('touchstart', handleTouchStart);
        panel.removeEventListener('touchmove', handleTouchMove);
        panel.removeEventListener('touchend', handleTouchEnd);
      }
    };
  }, [visible, onClose, side]);

  const panelClass = `dashboard-sidepanel dashboard-sidepanel--${side}${visible ? ' dashboard-sidepanel--visible' : ''}`;
  return (
    <aside ref={panelRef} className={panelClass}>
      {children}
    </aside>
  );
};

export default SidePanel;
