import { useState, useEffect } from 'react';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import InvestorList from './SidePanel/InvestorList/InvestorList';
import SidePanelButton from './SidePanel/SidePanelButton';
import SidePanel from './SidePanel/SidePanel';
import type { User } from '../User';
import { makeAuthenticatedRequest } from '../auth';

interface DashboardProps {
  onLogout: () => void;
}

export default function Dashboard({ onLogout }: DashboardProps) {
  const [user, setUser] = useState<User | null>(null);
  const [isLeftPanelOpen, setIsLeftPanelOpen] = useState(false);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const [loading, setLoading] = useState(true);

  const mainContent = `Welcome to the dashboard, ${user?.username}!`;

  useEffect(() => {
    const fetchUser = async () => {
      try {
        const response = await makeAuthenticatedRequest('/api/user');
        if (response.ok) {
          const data: User = await response.json();
          setUser(data);
        } else {
          onLogout();
        }
      } catch {
        onLogout();
      } finally {
        setLoading(false);
      }
    };
    fetchUser();
  }, [onLogout]);

  if (loading || !user) {
    return <div>Loading...</div>;
  }

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
          <h3>Settings</h3>
          <button
            onClick={onLogout}
            className='dashboard-logout-button'
          >
            Logout
          </button>
        </SidePanel>
      </div>
    </div>
  );
}
