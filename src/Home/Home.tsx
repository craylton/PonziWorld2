import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import './Home.css';
import { makeAuthenticatedRequest } from '../auth';
import type { Bank } from '../models/Bank';

interface HomeProps {
  onLogout: () => void;
}

export default function Home({ onLogout }: HomeProps) {
  const [banks, setBanks] = useState<Bank[]>([]);
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

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="home-container">
      <h1>My Banks</h1>
      <div className="bank-list">
        {banks.map((bank) => (
          <div
            key={bank.id}
            className="bank-card"
            onClick={() => handleBankClick(bank.id)}
          >
            <h2>{bank.bankName}</h2>
            <div className="bank-details">
              <div className="bank-stat">
                <span className="stat-label">Claimed Capital:</span>
                <span className="stat-value">${bank.claimedCapital.toLocaleString()}</span>
              </div>
              <div className="bank-stat">
                <span className="stat-label">Actual Capital:</span>
                <span className="stat-value">${bank.actualCapital.toLocaleString()}</span>
              </div>
              <div className="bank-stat">
                <span className="stat-label">Investors:</span>
                <span className="stat-value">{bank.investors.length}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
