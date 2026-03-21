import { useEffect, useState } from 'react';
import { makeAuthenticatedRequest } from '../../auth';
import { useLoadingContext } from '../../contexts/useLoadingContext';

interface TomorrowPanelProps {
  bankId: string;
  ponziFactor: number;
  onRefreshBank: () => Promise<void>;
}

const minPonziFactor = -0.05;
const maxPonziFactor = 0.15;

const sanitizePonziFactor = (value: number) => {
  if (!Number.isFinite(value)) {
    return 0;
  }

  return Math.max(minPonziFactor, Math.min(maxPonziFactor, value));
};

export default function TomorrowPanel({ bankId, ponziFactor, onRefreshBank }: TomorrowPanelProps) {
  const { showLoadingPopup } = useLoadingContext();
  const [ponziFactorInput, setPonziFactorInput] = useState<number>(() => sanitizePonziFactor(ponziFactor));

  useEffect(() => {
    setPonziFactorInput(sanitizePonziFactor(ponziFactor));
  }, [ponziFactor]);

  const formatPercent = (value: number) => {
    const percent = value * 100;
    const rounded = Math.round(percent * 10) / 10;
    if (Number.isInteger(rounded)) {
      return `${rounded.toFixed(0)}%`;
    }
    return `${rounded.toFixed(1)}%`;
  };

  const handleSavePonziFactor = async () => {
    if (
      !Number.isFinite(ponziFactorInput) ||
      ponziFactorInput < minPonziFactor ||
      ponziFactorInput > maxPonziFactor
    ) {
      showLoadingPopup('error', 'Please enter a valid Ponzi value.');
      return;
    }

    showLoadingPopup('loading', 'Saving Ponzi value...');

    try {
      const response = await makeAuthenticatedRequest('/api/bank/ponziFactor', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ bankId, ponziFactor: ponziFactorInput }),
      });

      if (response.ok) {
        await onRefreshBank();
        showLoadingPopup('success', 'Ponzi value saved.');
      } else {
        const errorData = await response.json();
        showLoadingPopup('error', `Failed to save Ponzi value: ${errorData.error || 'Unknown error'}.`);
      }
    } catch {
      showLoadingPopup('error', 'Failed to save Ponzi value: Network error.');
    }
  };

  return (
    <div>
      <label>
        Ponzi value:{' '}
        <span
          style={{
            display: 'inline-block',
            minWidth: '6ch',
            textAlign: 'right',
            fontVariantNumeric: 'tabular-nums',
          }}
        >
          {formatPercent(ponziFactorInput)}
        </span>
        <input
          type="range"
          min={minPonziFactor}
          max={maxPonziFactor}
          step={0.001}
          value={ponziFactorInput}
          onChange={(event) => setPonziFactorInput(sanitizePonziFactor(Number(event.target.value)))}
        />
      </label>
      <button
        type="button"
        onClick={handleSavePonziFactor}
        className="dashboard-settings-button"
      >
        Save Ponzi value
      </button>
    </div>
  );
}
