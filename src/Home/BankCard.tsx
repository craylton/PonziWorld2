import type { Bank } from '../models/Bank';
import { formatCurrency } from '../utils/currency';

interface BankCardProps {
  bank: Bank;
  onClick: (bankId: string) => void;
}

export default function BankCard({ bank, onClick }: BankCardProps) {
  return (
    <div
      className="bank-card"
      onClick={() => onClick(bank.id)}
    >
      <h2>{bank.bankName}</h2>
      <div className="bank-details">
        <div className="bank-stat">
          <span className="stat-label">Claimed Capital:</span>
          <span className="stat-value">{formatCurrency(bank.claimedCapital)}</span>
        </div>
        <div className="bank-stat">
          <span className="stat-label">Actual Capital:</span>
          <span className="stat-value">{formatCurrency(bank.actualCapital)}</span>
        </div>
        <div className="bank-stat">
          <span className="stat-label">Investors:</span>
          <span className="stat-value">{bank.investors.length}</span>
        </div>
      </div>
    </div>
  );
}
