import { useState, useEffect } from 'react';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import InvestorList from './SidePanel/InvestorList/InvestorList';
import SidePanelButton from './SidePanel/SidePanelButton';
import SidePanel from './SidePanel/SidePanel';
import type { User, Bank } from '../User';
import { makeAuthenticatedRequest, getUsernameFromToken } from '../auth';

interface DashboardProps {
  onLogout: () => void;
}

export default function Dashboard({ onLogout }: DashboardProps) {
  const [user, setUser] = useState<User | null>(null);
  const [bank, setBank] = useState<Bank | null>(null);
  const [isLeftPanelOpen, setIsLeftPanelOpen] = useState(false);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const [loading, setLoading] = useState(true);

  const mainContent = `Welcome to the dashboard, ${user?.username}!`;

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Get username from JWT and set user
        const username = getUsernameFromToken();
        if (!username) {
          onLogout();
          return;
        }
        setUser({ username });

        // Fetch bank data
        const bankResponse = await makeAuthenticatedRequest('/api/bank');
        if (!bankResponse.ok) {
          onLogout();
          return;
        }
        const bankData: Bank = await bankResponse.json();
        setBank(bankData);
      } catch {
        onLogout();
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, [onLogout]);

  if (loading || !user || !bank) {
    return <div>Loading...</div>;
  }

  return (
    <div className="dashboard-root">      <DashboardHeader
        bankName={bank.bankName}
        claimedCapital={bank.claimedCapital}
        actualCapital={bank.actualCapital}
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
