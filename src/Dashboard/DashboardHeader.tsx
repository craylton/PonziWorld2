import './Dashboard.css';

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
      <div className="dashboard-header__claimed">Claimed Capital: <span>{formatCurrency(claimedCapital)}</span></div>
      <div className="dashboard-header__actual">Actual Capital: <span>{formatCurrency(actualCapital)}</span></div>
    </header>
  );
}
