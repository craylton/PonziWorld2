import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './Home.css';
import { makeAuthenticatedRequest } from '../auth';
import type { Bank } from '../models/Bank';
import type { Player } from '../models/Player';
import SettingsButton from '../Dashboard/SidePanel/Settings/SettingsButton';
import SettingsPanel from '../Dashboard/SidePanel/Settings/SettingsPanel';
import LoadingProvider from '../contexts/LoadingContext';
import BankCard from './BankCard';

interface HomeProps {
  onLogout: () => void;
}

export default function Home({ onLogout }: HomeProps) {
  const [banks, setBanks] = useState<Bank[]>([]);
  const [player, setPlayer] = useState<Player | null>(null);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchBanks = async () => {
      try {
        const bankResponse = await makeAuthenticatedRequest('/api/banks');
        if (!bankResponse.ok) {
          onLogout();
          return;
        }
        const bankData: Bank[] = await bankResponse.json();
        setBanks(bankData);

        const playerResponse = await makeAuthenticatedRequest('/api/player');
        if (!playerResponse.ok) {
          onLogout();
          return;
        }
        const playerData: Player = await playerResponse.json();
        setPlayer(playerData);

        setIsLoading(false);
      } catch {
        onLogout();
      }
    };
    fetchBanks();
  }, [onLogout]);

  const handleBankClick = (bankId: string) => {
    navigate(`/bank/${bankId}`);
  };

  if (isLoading || !player) {
    return <div>Loading...</div>;
  }

  return (
    <LoadingProvider>
      <div className="home-container">
        <h1>My Banks</h1>
        <div className="bank-list">
          {banks.map((bank) => (
            <BankCard key={bank.id} bank={bank} onClick={handleBankClick} />
          ))}
        </div>
        <SettingsButton
          isRightPanelOpen={isRightPanelOpen}
          onClick={() => setIsRightPanelOpen(!isRightPanelOpen)}
        />
        <SettingsPanel
          visible={isRightPanelOpen}
          player={player}
          onLogout={onLogout}
          onClose={() => setIsRightPanelOpen(false)}
        />
      </div>
    </LoadingProvider>
  );
}
