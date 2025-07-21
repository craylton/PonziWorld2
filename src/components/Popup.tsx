import { useEffect } from 'react';
import type { ReactNode } from 'react';
import './Popup.css';

interface PopupProps {
  isOpen: boolean;
  title: string;
  onClose: () => void;
  children: ReactNode;
  footer?: ReactNode;
  className?: string;
  zIndex?: number;
  canCloseOnOverlayClick?: boolean;
  showCloseButton?: boolean;
}

export default function Popup({
  isOpen,
  title,
  onClose,
  children,
  footer,
  className = '',
  zIndex = 2000,
  canCloseOnOverlayClick = true,
  showCloseButton = true
}: PopupProps) {
  // Prevent background scrolling when open
  useEffect(() => {
    if (isOpen) document.body.style.overflow = 'hidden';
    return () => { document.body.style.overflow = 'unset'; };
  }, [isOpen]);

  if (!isOpen) return null;

  const handleOverlayClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (canCloseOnOverlayClick && e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div
      className="popup-overlay"
      style={{ zIndex }}
      onClick={handleOverlayClick}
      role="dialog"
      aria-modal="true"
      aria-labelledby="popup-title"
    >
      <div className={`popup ${className}`}>
        <div className="popup__header">
          <h2 id="popup-title" className="popup__title">{title}</h2>
          {showCloseButton && (
            <button
              className="popup__close-button"
              onClick={onClose}
              aria-label="Close popup"
            >
              ×
            </button>
          )}
        </div>
        <div className="popup__content">
          {children}
        </div>
        {footer && (
          <div className="popup__footer">
            {footer}
          </div>
        )}
      </div>
    </div>
  );
}
