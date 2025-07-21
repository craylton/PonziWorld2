import { useState, useEffect } from 'react';
import '../CapitalPopupStyles.css';
import Popup from '../../components/Popup';

interface TransactionPopupProps {
    isOpen: boolean;
    onClose: () => void;
    assetType: string;
    transactionType: 'buy' | 'sell';
    currentHoldings?: number;
    onConfirm: (amount: number) => void;
}

export default function TransactionPopup({
    isOpen,
    onClose,
    assetType,
    transactionType,
    currentHoldings = 0,
    onConfirm
}: TransactionPopupProps) {
    const [amount, setAmount] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [sellAll, setSellAll] = useState<boolean>(false);

    const title = transactionType === 'buy' ? `Buy ${assetType}` : `Sell ${assetType}`;

    // Reset form when popup opens/closes
    useEffect(() => {
        if (isOpen) {
            setAmount('');
            setError('');
            setSellAll(false);
        }
    }, [isOpen]);

    const validateAmount = (value: string): string => {
        const numValue = parseFloat(value);

        if (isNaN(numValue) || numValue <= 0) {
            return 'Amount must be a positive number';
        }

        if (transactionType === 'sell' && numValue > currentHoldings) {
            return `Cannot sell more than your current holdings (${currentHoldings})`;
        }

        return '';
    };

    const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setAmount(value);

        // Check if the entered amount equals current holdings (for sell all functionality)
        if (transactionType === 'sell' && value && parseFloat(value) === currentHoldings) {
            setSellAll(true);
        } else {
            setSellAll(false);
        }

        if (value) {
            const validationError = validateAmount(value);
            setError(validationError);
        } else {
            setError('');
        }
    };

    const handleSellAllChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const checked = e.target.checked;
        setSellAll(checked);

        if (checked) {
            setAmount(currentHoldings.toString());
            setError('');
        } else {
            setAmount('');
        }
    };

    const handleConfirm = () => {
        const validationError = validateAmount(amount);
        if (validationError) {
            setError(validationError);
            return;
        }

        const numAmount = parseFloat(amount);
        onConfirm(numAmount);
        onClose();
    };

    const footer = (
        <>
            <button
                className="popup__button popup__button--secondary"
                onClick={onClose}
            >
                Cancel
            </button>
            <button
                className="popup__button popup__button--confirm"
                onClick={handleConfirm}
                disabled={!amount || !!error}
            >
                {transactionType === 'buy' ? 'Buy' : 'Sell'}
            </button>
        </>
    );

    return (
        <Popup
            isOpen={isOpen}
            title={title}
            onClose={onClose}
            footer={footer}
            zIndex={2001}
            className="transaction-popup"
        >
            <div className="transaction-popup__input-group">
                <label htmlFor="amount-input" className="transaction-popup__label">
                    Amount:
                </label>
                <input
                    id="amount-input"
                    type="number"
                    value={amount}
                    onChange={handleAmountChange}
                    placeholder="Enter amount"
                    className="transaction-popup__input"
                    min="0"
                    step="0.01"
                    autoFocus
                />
                {error && (
                    <div className="transaction-popup__error">
                        {error}
                    </div>
                )}
            </div>
            {transactionType === 'sell' && (
                <div className="transaction-popup__checkbox-group">
                    <label className="transaction-popup__checkbox-label">
                        <input
                            type="checkbox"
                            checked={sellAll}
                            onChange={handleSellAllChange}
                            className="transaction-popup__checkbox"
                        />
                        Sell all
                    </label>
                </div>
            )}
        </Popup>
    );
}
