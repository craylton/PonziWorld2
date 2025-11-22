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
import CreateBankPopup from './CreateBankPopup';
import LoadingPopup from '../Dashboard/Assets/LoadingPopup';
import PageHeader from '../components/PageHeader';

interface HomeProps {
  onLogout: () => void;
}

export default function Home({ onLogout }: HomeProps) {
  const [banks, setBanks] = useState<Bank[]>([]);
  const [player, setPlayer] = useState<Player | null>(null);
  const [isRightPanelOpen, setIsRightPanelOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [isCreateBankPopupOpen, setIsCreateBankPopupOpen] = useState(false);
  const [loadingPopupState, setLoadingPopupState] = useState<{
    isOpen: boolean;
    status: 'loading' | 'success' | 'error';
    message?: string;
  }>({ isOpen: false, status: 'loading' });
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

  const handleCreateBank = async (bankName: string) => {
    setIsCreateBankPopupOpen(false);
    setLoadingPopupState({ isOpen: true, status: 'loading', message: 'Creating bank...' });

    try {
      const response = await makeAuthenticatedRequest('/api/bank', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ bankName }),
      });

      if (response.ok) {
        setLoadingPopupState({ isOpen: true, status: 'success', message: 'Bank created successfully' });
        // Refresh the banks list
        const bankResponse = await makeAuthenticatedRequest('/api/banks');
        if (bankResponse.ok) {
          const bankData: Bank[] = await bankResponse.json();
          setBanks(bankData);
        }
      } else {
        const errorData = await response.json();
        const errorMessage = errorData.error || 'Failed to create bank';
        setLoadingPopupState({ isOpen: true, status: 'error', message: errorMessage });
      }
    } catch {
      setLoadingPopupState({ isOpen: true, status: 'error', message: 'Network error occurred' });
    }
  };

  const handleCloseLoadingPopup = () => {
    setLoadingPopupState({ isOpen: false, status: 'loading' });
  };

  if (isLoading || !player) {
    return <div>Loading...</div>;
  }

  return (
    <LoadingProvider>
      <div className="home-container">
        <PageHeader title="My Banks" />
        <div className="home-content">
          <div className="bank-list">
            {banks.map((bank) => (
              <BankCard key={bank.id} bank={bank} onClick={handleBankClick} />
            ))}
          </div>
          <button
            className="create-bank-button"
            onClick={() => setIsCreateBankPopupOpen(true)}
            disabled={banks.length >= 3}
          >
            Create New Bank
          </button>
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
        <CreateBankPopup
          isOpen={isCreateBankPopupOpen}
          onClose={() => setIsCreateBankPopupOpen(false)}
          onConfirm={handleCreateBank}
        />
        <LoadingPopup
          isOpen={loadingPopupState.isOpen}
          onClose={handleCloseLoadingPopup}
          status={loadingPopupState.status}
          message={loadingPopupState.message}
        />
      </div>
    </LoadingProvider>
  );
}
