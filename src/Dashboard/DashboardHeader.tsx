import './DashboardHeader.css';
import type { PerformanceHistory } from '../User';

interface DashboardHeaderProps {
  bankName: string;
  claimedCapital: number;
  actualCapital: number;
  performanceHistory: PerformanceHistory | null;
}

function formatCurrency(amount: number) {
  return amount.toLocaleString(undefined, { style: 'currency', currency: 'GBP', maximumFractionDigits: 2 });
}

export default function DashboardHeader({ bankName, claimedCapital, actualCapital, performanceHistory }: DashboardHeaderProps) {
  // Performance history is available here for future use (charts, trends, etc.)
  // Currently not displayed but loaded for future features
  void performanceHistory; // Mark as intentionally unused
  
  return (
    <header className="dashboard-header">
      <div className="dashboard-header__bank">{bankName}</div>
      <div className="dashboard-header__capitals">
        <div className="dashboard-header__capital">
          <span className="dashboard-header__capital-label">Claimed Capital</span>
          <span className="dashboard-header__capital-value">{formatCurrency(claimedCapital)}</span>
        </div>
        <div className="dashboard-header__capital">
          <span className="dashboard-header__capital-label">Actual Capital</span>
          <span className="dashboard-header__capital-value">{formatCurrency(actualCapital)}</span>
        </div>
      </div>
    </header>
  );
}
