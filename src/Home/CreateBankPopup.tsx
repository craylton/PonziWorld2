import { useState } from 'react';
import Popup from '../components/Popup';
import './CreateBankPopup.css';

interface CreateBankPopupProps {
  isOpen: boolean;
  onClose: () => void;
  onConfirm: (bankName: string) => void;
}

export default function CreateBankPopup({
  isOpen,
  onClose,
  onConfirm
}: CreateBankPopupProps) {
  const [bankName, setBankName] = useState('');

  const handleConfirm = () => {
    if (bankName.trim()) {
      onConfirm(bankName.trim());
      setBankName('');
    }
  };

  const handleClose = () => {
    setBankName('');
    onClose();
  };

  const footer = (
    <>
      <button
        className="popup__button popup__button--cancel"
        onClick={handleClose}
      >
        Cancel
      </button>
      <button
        className="popup__button popup__button--confirm"
        onClick={handleConfirm}
        disabled={!bankName.trim()}
      >
        Confirm
      </button>
    </>
  );

  return (
    <Popup
      isOpen={isOpen}
      title="Create New Bank"
      onClose={handleClose}
      footer={footer}
      className="create-bank-popup"
    >
      <div className="create-bank-popup__content">
        <label htmlFor="bank-name" className="create-bank-popup__label">
          Bank Name
        </label>
        <input
          id="bank-name"
          type="text"
          className="create-bank-popup__input"
          value={bankName}
          onChange={(e) => setBankName(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter' && bankName.trim()) {
              handleConfirm();
            }
          }}
          placeholder="Enter bank name"
          autoFocus
        />
      </div>
    </Popup>
  );
}
