import '../CapitalPopup.css';
import './LoadingPopup.css';

interface LoadingPopupProps {
    isOpen: boolean;
    onClose: () => void;
    status: 'loading' | 'success' | 'error';
    message?: string;
}

export default function LoadingPopup({
    isOpen,
    onClose,
    status,
    message
}: LoadingPopupProps) {
    if (!isOpen) return null;

    const getStatusMessage = () => {
        switch (status) {
            case 'loading':
                return 'Loading...';
            case 'success':
                return message || 'Request succeeded';
            case 'error':
                return message || 'Request failed';
            default:
                return 'Loading...';
        }
    };

    const getStatusClass = () => {
        switch (status) {
            case 'success':
                return 'loading-popup__message--success';
            case 'error':
                return 'loading-popup__message--error';
            default:
                return '';
        }
    };

    const canClose = status === 'success' || status === 'error';

    const onClickOutside = (e: React.MouseEvent<HTMLDivElement>) => {
        if (canClose && e.target === e.currentTarget) {
            onClose();
        }
    }

    return (
        <div
            className="capital-popup-overlay"
            onClick={onClickOutside}
            role="dialog"
            aria-modal="true"
            aria-labelledby="loading-popup-title"
            style={{ zIndex: 2500 }} // Higher z-index to appear above other popups
        >
            <div className="capital-popup loading-popup">
                <div className="capital-popup__header">
                    <h2 id="loading-popup-title" className="capital-popup__title">
                        Transaction Status
                    </h2>
                    {canClose && (
                        <button
                            className="capital-popup__close-button"
                            onClick={onClose}
                            aria-label="Close popup"
                        >
                            ×
                        </button>
                    )}
                </div>
                <div className="capital-popup__content">
                    <div className={`loading-popup__message ${getStatusClass()}`}>
                        {status === 'loading' && (
                            <div className="loading-popup__spinner" aria-label="Loading..." />
                        )}
                        <span>{getStatusMessage()}</span>
                    </div>
                </div>
                {canClose && (
                    <div className="capital-popup__footer">
                        <button
                            className="capital-popup__confirm-button"
                            onClick={onClose}
                        >
                            OK
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}
