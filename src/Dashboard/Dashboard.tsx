import { useState, useEffect, useCallback } from 'react';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import InvestorsButton from './SidePanel/Investors/InvestorsButton';
import SettingsButton from './SidePanel/Settings/SettingsButton';
import InvestorsPanel from './SidePanel/Investors/InvestorsPanel';
import SettingsPanel from './SidePanel/Settings/SettingsPanel';
import AssetSection from './Assets/AssetSection';
import { makeAuthenticatedRequest } from '../auth';
import { BankProvider } from '../contexts/BankContext';
import type { Bank } from '../models/Bank';
import type { HistoricalPerformance } from '../models/HistoricalPerformance';
import type { Player } from '../models/User';

interface DashboardProps {
  onLogout: () => void;
}

export default function Dashboard({ onLogout }: DashboardProps) {
  const [bank, setBank] = useState<Bank | null>(null);
  const [player, setPlayer] = useState<Player | null>(null);
  const [historicalPerformance, setHistoricalPerformance] = useState<HistoricalPerformance | null>(null);
  const [currentDay, setCurrentDay] = useState<number | null>(null);
  const [isLeftPanelOpen, setIsLeftPanelOpen] = useState(false);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const [isInitialDataLoading, setIsInitialDataLoading] = useState(true);
  const [isHistoryLoading, setIsHistoryLoading] = useState(true);

  const fetchBankData = useCallback(async (currentBankId?: string) => {
    try {
      // Fetch bank data
      const bankResponse = await makeAuthenticatedRequest('/api/banks');
      if (!bankResponse.ok) {
        onLogout();
        return null;
      }
      const bankData: Bank[] = await bankResponse.json();
      // For now, we'll just use the first bank
      const firstBank = bankData[0];
      setBank(firstBank);

      // If this is an update (not initial load), refresh historical performance too
      if (currentBankId) {
        const historyResponse = await makeAuthenticatedRequest(`/api/historicalPerformance/ownbank/${firstBank.id}`);
        if (historyResponse.ok) {
          const historyData: HistoricalPerformance = await historyResponse.json();
          setHistoricalPerformance(historyData);
        }
      }

      return firstBank;
    } catch {
      onLogout();
      return null;
    }
  }, [onLogout]);

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch current day
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
        const firstBank = await fetchBankData();
        if (!firstBank) {
          return;
        }

        // All essential data pieces loaded
        setIsInitialDataLoading(false);

        // Fetch performance history (non-essential, can load separately)
        const historyResponse = await makeAuthenticatedRequest(`/api/historicalPerformance/ownbank/${firstBank.id}`);
        if (historyResponse.ok) {
          const historyData: HistoricalPerformance = await historyResponse.json();
          setHistoricalPerformance(historyData);
        }
      } catch {
        onLogout();
      } finally {
        setIsHistoryLoading(false);
      }
    };
    fetchData();
  }, [onLogout, fetchBankData]);

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
        historicalPerformance={historicalPerformance}
        isHistoryLoading={isHistoryLoading}
      />
      <div className="dashboard-layout">
        <InvestorsPanel 
          visible={isLeftPanelOpen} 
          onClose={() => setIsLeftPanelOpen(false)}
        />
        <main className="dashboard-main">
          <InvestorsButton
            isLeftPanelOpen={isLeftPanelOpen}
            onClick={() => setIsLeftPanelOpen(!isLeftPanelOpen)}
          />
          <BankProvider bankId={bank.id}>
            <AssetSection 
              availableAssets={bank.availableAssets} 
              onRefreshBank={async () => {
                await fetchBankData(bank.id);
              }}
            />
          </BankProvider>

          <SettingsButton
            isRightPanelOpen={isRightPanelOpen}
            onClick={() => setIsRightPanelOpen(!isRightPanelOpen)}
          />
        </main>
        <SettingsPanel
          visible={isRightPanelOpen}
          player={player}
          onLogout={onLogout}
          onClose={() => setIsRightPanelOpen(false)}
        />
      </div>
    </div>
  );
}
