import { useState } from 'react';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import InvestorList from './SidePanel/InvestorList/InvestorList';
import SidePanelButton from './SidePanel/SidePanelButton';
import SidePanel from './SidePanel/SidePanel';

interface User {
  id: string;
  username: string;
  bankName: string;
  claimedCapital: number;
  actualCapital: number;
}

interface DashboardProps {
  user: User;
}

export default function Dashboard({ user }: DashboardProps) {
  const [isLeftPanelOpen, setIsLeftPanelOpen] = useState(false);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const mainContent = `Welcome to the dashboard, ${user.username}!`;

  return (
    <div className="dashboard-root">
      <DashboardHeader
        bankName={user.bankName}
        claimedCapital={user.claimedCapital}
        actualCapital={user.actualCapital}
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
