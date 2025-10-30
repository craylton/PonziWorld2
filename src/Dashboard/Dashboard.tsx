import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import InvestorsButton from './SidePanel/Investors/InvestorsButton';
import SettingsButton from './SidePanel/Settings/SettingsButton';
import InvestorsPanel from './SidePanel/Investors/InvestorsPanel';
import SettingsPanel from './SidePanel/Settings/SettingsPanel';
import AssetSection from './Assets/AssetSection';
import { makeAuthenticatedRequest } from '../auth';
import { BankProvider } from '../contexts/BankContext';
import LoadingProvider from '../contexts/LoadingContext';
import type { Bank } from '../models/Bank';
import type { Player } from '../models/Player';

interface DashboardProps {
  onLogout: () => void;
}

export default function Dashboard({ onLogout }: DashboardProps) {
  const { bankId } = useParams<{ bankId: string }>();
  const navigate = useNavigate();
  const [bank, setBank] = useState<Bank | null>(null);
  const [player, setPlayer] = useState<Player | null>(null);
  const [currentDay, setCurrentDay] = useState<number | null>(null);
  const [isLeftPanelOpen, setIsLeftPanelOpen] = useState(false);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const [isInitialDataLoading, setIsInitialDataLoading] = useState(true);

  const fetchBankData = useCallback(async () => {
    if (!bankId) {
      navigate('/');
      return null;
    }

    try {
      // Fetch bank data
      const bankResponse = await makeAuthenticatedRequest('/api/banks');
      if (!bankResponse.ok) {
        onLogout();
        return null;
      }
      const bankData: Bank[] = await bankResponse.json();
      
      // Find the specific bank by ID
      const specificBank = bankData.find(b => b.id === bankId);
      if (!specificBank) {
        // Bank not found, redirect to home
        navigate('/');
        return null;
      }
      
      setBank(specificBank);
      return specificBank;
    } catch {
      onLogout();
      return null;
    }
  }, [onLogout, bankId, navigate]);

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
      } catch {
        onLogout();
      }
    };
    fetchData();
  }, [onLogout, fetchBankData]);

  if (isInitialDataLoading || !bank || !player || currentDay === null) {
    return <div>Loading...</div>;
  }

  return (
    <LoadingProvider>
      <div className="dashboard-root">
        <DashboardHeader
          currentDay={currentDay}
          bankName={bank.bankName}
          claimedCapital={bank.claimedCapital}
          actualCapital={bank.actualCapital}
          bankId={bank.id}
        />
        <div className="dashboard-layout">
          <InvestorsPanel
            visible={isLeftPanelOpen}
            onClose={() => setIsLeftPanelOpen(false)}
            investors={bank.investors}
          />
          <main className="dashboard-main">
            <InvestorsButton
              isLeftPanelOpen={isLeftPanelOpen}
              onClick={() => setIsLeftPanelOpen(!isLeftPanelOpen)}
            />
            <BankProvider bankId={bank.id}>
              <AssetSection
                availableAssets={bank.availableAssets}
                onRefreshBank={async () => { await fetchBankData(); }}
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
    </LoadingProvider>
  );
}
