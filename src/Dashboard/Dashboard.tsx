import { useState } from 'react';
import SidePanel from './SidePanel';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import HamburgerButton from './HamburgerButton';
import InvestorList from './InvestorList';

interface DashboardProps {
  username: string;
}

// TODO: Replace with real data from backend
const DUMMY_BANK_NAME = 'Ponzi National Bank';
const DUMMY_CLAIMED_CAPITAL = 1000000;
const DUMMY_ACTUAL_CAPITAL = 250000;

export default function Dashboard({ username }: DashboardProps) {
  const [showLeftPanel, setShowLeftPanel] = useState(false);
  const [showRightPanel, setShowRightPanel] = useState(false);
  const mainContent = `Welcome to the dashboard, ${username}!`;

  return (
    <div className="dashboard-root">
      <DashboardHeader
        bankName={DUMMY_BANK_NAME}
        claimedCapital={DUMMY_CLAIMED_CAPITAL}
        actualCapital={DUMMY_ACTUAL_CAPITAL}
      />
      <div className="dashboard-layout">        {/* Left Side Panel */}
        <SidePanel side="left" visible={showLeftPanel}>
          <HamburgerButton
            isOpen={showLeftPanel}
            onClick={() => setShowLeftPanel((v) => !v)}
            ariaLabel={showLeftPanel ? 'Close left panel' : 'Open left panel'}
            className="dashboard-hamburger--panel"
          />
          <InvestorList />
        </SidePanel>
        <main className="dashboard-main">
          {/* Hamburger for left panel (mobile only, floats over main) */}
          <HamburgerButton
            isOpen={false}
            onClick={() => setShowLeftPanel(true)}
            ariaLabel="Open left panel"
            className={`dashboard-hamburger--float dashboard-hamburger--left${showLeftPanel ? ' dashboard-hamburger--hidden' : ''}`}
          />
          {mainContent}
          {/* Hamburger for right panel (mobile only, floats over main) */}
          <HamburgerButton
            isOpen={false}
            onClick={() => setShowRightPanel(true)}
            ariaLabel="Open right panel"
            className={`dashboard-hamburger--float dashboard-hamburger--right${showRightPanel ? ' dashboard-hamburger--hidden' : ''}`}
          />
        </main>
        {/* Right Side Panel */}
        <SidePanel side="right" visible={showRightPanel}>
          <HamburgerButton
            isOpen={showRightPanel}
            onClick={() => setShowRightPanel((v) => !v)}
            ariaLabel={showRightPanel ? 'Close right panel' : 'Open right panel'}
            className="dashboard-hamburger--panel"
          />
          Right Panel
        </SidePanel>
      </div>
    </div>
  );
}
