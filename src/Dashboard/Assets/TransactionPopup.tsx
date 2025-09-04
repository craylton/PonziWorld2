import { useState, useEffect } from 'react';
import { parseMoney, compareMoney, isValidMoney } from '../../utils/money';
import Popup from '../../components/Popup';

interface TransactionPopupProps {
    isOpen: boolean;
    onClose: () => void;
    assetName: string;
    transactionType: 'buy' | 'sell';
    currentHoldings?: string; // Now string for arbitrary precision
    maxBuyAmount?: string;   // Now string for arbitrary precision
    onConfirm: (amount: string) => void; // Now string for arbitrary precision
}

export default function TransactionPopup({
    isOpen,
    onClose,
    assetName,
    transactionType,
    currentHoldings = '0',  // Now string default
    maxBuyAmount = '0',     // Now string default
    onConfirm
}: TransactionPopupProps) {
    const [amount, setAmount] = useState<string>('');
    const [error, setError] = useState<string>('');
    const [sellAll, setSellAll] = useState<boolean>(false);

    const title = transactionType === 'buy' ? `Buy ${assetName}` : `Sell ${assetName}`;

    // Reset form when popup opens/closes
    useEffect(() => {
        if (isOpen) {
            setAmount('');
            setError('');
            setSellAll(false);
        }
    }, [isOpen]);

    const validateAmount = (value: string): string => {
        if (!isValidMoney(value)) {
            return 'Amount must be a valid number';
        }

        const amount = parseMoney(value);
        
        if (amount.lte(0)) {
            return 'Amount must be a positive number';
        }

        if (transactionType === 'sell' && compareMoney(amount, parseMoney(currentHoldings)) > 0) {
            return `Cannot sell more than your current holdings (${currentHoldings})`;
        }

        if (transactionType === 'buy' && compareMoney(amount, parseMoney(maxBuyAmount)) > 0) {
            return `Cannot buy more than your available cash (${maxBuyAmount})`;
        }

        return '';
    };

    const handleAmountChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value = e.target.value;
        setAmount(value);

        // Check if the entered amount equals current holdings (for sell all functionality)
        if (transactionType === 'sell' && value && isValidMoney(value)) {
            const enteredAmount = parseMoney(value);
            const holdingsAmount = parseMoney(currentHoldings);
            if (compareMoney(enteredAmount, holdingsAmount) === 0) {
                setSellAll(true);
            } else {
                setSellAll(false);
            }
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

        onConfirm(amount); // Pass string directly for arbitrary precision
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
            <div className="popup__input-group">
                <label htmlFor="amount-input" className="popup__label">
                    Amount:
                </label>
                <input
                    id="amount-input"
                    type="number"
                    value={amount}
                    onChange={handleAmountChange}
                    placeholder="Enter amount"
                    className="popup__input"
                    min="0"
                    step="0.01"
                    autoFocus
                />
                {error && (
                    <div className="popup__error">
                        {error}
                    </div>
                )}
            </div>
            {transactionType === 'sell' && (
                <div className="popup__checkbox-group">
                    <label className="popup__checkbox-label">
                        <input
                            type="checkbox"
                            checked={sellAll}
                            onChange={handleSellAllChange}
                            className="popup__checkbox"
                        />
                        Sell all
                    </label>
                </div>
            )}
        </Popup>
    );
}
