import { useState, useEffect } from 'react';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import InvestorList from './SidePanel/InvestorList/InvestorList';
import SidePanelButton from './SidePanel/SidePanelButton';
import SidePanel from './SidePanel/SidePanel';
import AssetSection from './AssetSection';
import { makeAuthenticatedRequest } from '../auth';
import type { Bank } from '../models/Bank';
import type { PerformanceHistory } from '../models/PerformanceHistory';
import type { Player } from '../models/User';

interface DashboardProps {
  onLogout: () => void;
}

export default function Dashboard({ onLogout }: DashboardProps) {
  const [bank, setBank] = useState<Bank | null>(null);
  const [player, setPlayer] = useState<Player | null>(null);
  const [performanceHistory, setPerformanceHistory] = useState<PerformanceHistory | null>(null);
  const [currentDay, setCurrentDay] = useState<number | null>(null);
  const [isLeftPanelOpen, setIsLeftPanelOpen] = useState(false);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const [isInitialDataLoading, setIsInitialDataLoading] = useState(true);
  const [isHistoryLoading, setIsHistoryLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch current day (non-authenticated)
        const currentDayResponse = await fetch('/api/currentDay');
        if (!currentDayResponse.ok) {
          onLogout();
          return;
        }
        const currentDayData: { currentDay: number } = await currentDayResponse.json();
        setCurrentDay(currentDayData.currentDay);

        // Fetch player data
        const playerResponse = await makeAuthenticatedRequest('/api/player');
        if (!playerResponse.ok) {
          onLogout();
          return;
        }
        const playerData: Player = await playerResponse.json();
        setPlayer(playerData);

        // Fetch bank data
        const bankResponse = await makeAuthenticatedRequest('/api/bank');
        if (!bankResponse.ok) {
          onLogout();
          return;
        }
        const bankData: Bank = await bankResponse.json();
        setBank(bankData);

        // All essential data pieces loaded
        setIsInitialDataLoading(false);

        // Fetch performance history (non-essential, can load separately)
        const historyResponse = await makeAuthenticatedRequest(`/api/performanceHistory/ownbank/${bankData.id}`);
        if (historyResponse.ok) {
          const historyData: PerformanceHistory = await historyResponse.json();
          setPerformanceHistory(historyData);
        }
      } catch {
        onLogout();
      } finally {
        setIsHistoryLoading(false);
      }
    };
    fetchData();
  }, [onLogout]);

  const handleAdvanceDay = async () => {
    try {
      const response = await makeAuthenticatedRequest('/api/nextDay', {
        method: 'POST',
      });
      
      if (response.ok) {
        // Optionally refresh the page or show a success message
        window.location.reload();
      } else {
        const errorData = await response.json();
        alert(`Failed to advance day: ${errorData.error || 'Unknown error'}`);
      }
    } catch {
      alert('Failed to advance day: Network error');
    }
  };

  if (isInitialDataLoading || !bank || !player || currentDay === null) {
    return <div>Loading...</div>;
  }

  return (
    <div className="dashboard-root">
      <DashboardHeader
        currentDay={currentDay}
        bankName={bank.bankName}
        claimedCapital={bank.claimedCapital}
        actualCapital={bank.actualCapital}
        performanceHistory={performanceHistory}
        isHistoryLoading={isHistoryLoading}
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
          <AssetSection bankAssets={bank.assets} />
          
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
            className='dashboard-settings-button'
          >
            Logout
          </button>
          {player.isAdmin && (
            <div className="dashboard-admin-section">
              <p>Admin only</p>
              <button
                onClick={handleAdvanceDay}
                className="dashboard-settings-button"
              >
                Advance to next day
              </button>
            </div>
          )}
        </SidePanel>
      </div>
    </div>
  );
}
