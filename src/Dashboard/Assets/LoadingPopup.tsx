import './LoadingPopup.css';
import Popup from '../../components/Popup';

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

    const footer = canClose ? (
        <button
            className="popup__button popup__button--confirm"
            onClick={onClose}
        >
            OK
        </button>
    ) : undefined;

    return (
        <Popup
            isOpen={isOpen}
            title="Transaction Status"
            onClose={onClose}
            footer={footer}
            canCloseOnOverlayClick={canClose}
            showCloseButton={canClose}
            className="loading-popup"
        >
            <div className={`loading-popup__message ${getStatusClass()}`}>
                {status === 'loading' && (
                    <div className="loading-popup__spinner" aria-label="Loading..." />
                )}
                <span>{getStatusMessage()}</span>
            </div>
        </Popup>
    );
}
