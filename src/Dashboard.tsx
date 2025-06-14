import DashboardHeader from './DashboardHeader';
import './Dashboard.css';

interface DashboardProps {
  username: string;
}

// TODO: Replace with real data from backend
const DUMMY_BANK_NAME = 'Ponzi National Bank';
const DUMMY_CLAIMED_CAPITAL = 1000000;
const DUMMY_ACTUAL_CAPITAL = 250000;

export default function Dashboard({ username }: DashboardProps) {
    const mainContent = `Welcome to the dashboard, ${username}!`;
  return (
    <div className="dashboard-root">
      <DashboardHeader
        bankName={DUMMY_BANK_NAME}
        claimedCapital={DUMMY_CLAIMED_CAPITAL}
        actualCapital={DUMMY_ACTUAL_CAPITAL}
      />
      {/* Layout: left panel, main, right panel */}
      <div className="dashboard-layout">
        <aside className="dashboard-panel dashboard-panel--left">Left Panel</aside>
        <main className="dashboard-main">{mainContent}</main>
        <aside className="dashboard-panel dashboard-panel--right">Right Panel</aside>
      </div>
    </div>
  );
}
