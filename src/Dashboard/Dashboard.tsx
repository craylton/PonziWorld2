import { useState } from 'react';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import InvestorList from './SidePanel/InvestorList/InvestorList';
import SidePanelButton from './SidePanel/SidePanelButton';
import SidePanel from './SidePanel/SidePanel';

interface DashboardProps {
  username: string;
}

// TODO: Replace with real data from backend
const DUMMY_BANK_NAME = 'Ponzi National Bank';
const DUMMY_CLAIMED_CAPITAL = 1000000;
const DUMMY_ACTUAL_CAPITAL = 250000;

export default function Dashboard({ username }: DashboardProps) {
  const [isLeftPanelOpen, setIsLeftPanelOpen] = useState(false);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const mainContent = `Welcome to the dashboard, ${username}!`;

  return (
    <div className="dashboard-root">
      <DashboardHeader
        bankName={DUMMY_BANK_NAME}
        claimedCapital={DUMMY_CLAIMED_CAPITAL}
        actualCapital={DUMMY_ACTUAL_CAPITAL}
      />
      <div className="dashboard-layout">
        <SidePanel side="left" visible={isLeftPanelOpen}>
          <InvestorList />
        </SidePanel>
        <main className="dashboard-main">
          <SidePanelButton
            iconType="hamburger"
            shouldAllowClose={isLeftPanelOpen}
            onClick={() => setIsLeftPanelOpen(!isLeftPanelOpen)}
            ariaLabel="Open left panel"
            className={`dashboard-sidepanel-button--left`}
          />
          {mainContent}
            <SidePanelButton
              iconType="cog"
              shouldAllowClose={isRightPanelOpen}
              onClick={() => setIsRightPanelOpen(!isRightPanelOpen)}
              ariaLabel="Open settings panel"
              className={`dashboard-sidepanel-button--right`}
            />
        </main>
        <SidePanel side="right" visible={isRightPanelOpen}>
          {/* TODO: Add settings content here */}
          <></>
        </SidePanel>
      </div>
    </div>
  );
}
