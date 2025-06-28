import { useState, useEffect } from 'react';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import InvestorList from './SidePanel/InvestorList/InvestorList';
import SidePanelButton from './SidePanel/SidePanelButton';
import SidePanel from './SidePanel/SidePanel';
import AssetList from './AssetList/AssetList';
import type { Bank, PerformanceHistory } from '../User';
import { makeAuthenticatedRequest } from '../auth';

interface DashboardProps {
  onLogout: () => void;
}

export default function Dashboard({ onLogout }: DashboardProps) {
  const [bank, setBank] = useState<Bank | null>(null);
  const [performanceHistory, setPerformanceHistory] = useState<PerformanceHistory | null>(null);
  const [isLeftPanelOpen, setIsLeftPanelOpen] = useState(false);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const [bankLoading, setBankLoading] = useState(true);
  const [historyLoading, setHistoryLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch bank data
        const bankResponse = await makeAuthenticatedRequest('/api/bank');
        if (!bankResponse.ok) {
          onLogout();
          return;
        }
        const bankData: Bank = await bankResponse.json();
        setBank(bankData);
        setBankLoading(false);

        // Fetch performance history
        const historyResponse = await makeAuthenticatedRequest(`/api/performanceHistory/ownbank/${bankData.id}`);
        if (historyResponse.ok) {
          const historyData: PerformanceHistory = await historyResponse.json();
          setPerformanceHistory(historyData);
        }
      } catch {
        onLogout();
      } finally {
        setHistoryLoading(false);
      }
    };
    fetchData();
  }, [onLogout]);

  if (bankLoading || !bank) {
    return <div>Loading...</div>;
  }

  return (
    <div className="dashboard-root">
      <DashboardHeader
        bankName={bank.bankName}
        claimedCapital={bank.claimedCapital}
        actualCapital={bank.actualCapital}
        performanceHistory={performanceHistory}
        historyLoading={historyLoading}
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
          <AssetList assets={bank.assets} />
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
