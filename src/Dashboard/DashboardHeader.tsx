import './DashboardHeader.css';

interface DashboardHeaderProps {
  bankName: string;
  claimedCapital: number;
  actualCapital: number;
}

function formatCurrency(amount: number) {
  return amount.toLocaleString(undefined, { style: 'currency', currency: 'GBP', maximumFractionDigits: 2 });
}

export default function DashboardHeader({ bankName, claimedCapital, actualCapital }: DashboardHeaderProps) {
  return (
    <header className="dashboard-header">
      <div className="dashboard-header__bank">{bankName}</div>
      <div className="dashboard-header__capitals">
        <div className="dashboard-header__capital">
          <span className="dashboard-header__capital-label">Claimed Capital</span>
          <span className="dashboard-header__capital-value">{formatCurrency(claimedCapital)}</span>
        </div>
        <div className="dashboard-header__separator" aria-hidden="true"></div>
        <div className="dashboard-header__capital">
          <span className="dashboard-header__capital-label">Actual Capital</span>
          <span className="dashboard-header__capital-value">{formatCurrency(actualCapital)}</span>
        </div>
      </div>
    </header>
  );
}
