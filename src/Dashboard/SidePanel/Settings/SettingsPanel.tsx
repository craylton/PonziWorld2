import SidePanel from '../SidePanel';
import type { Player } from '../../../models/Player';
import { makeAuthenticatedRequest } from '../../../auth';
import { useLoadingContext } from '../../../contexts/useLoadingContext';

interface SettingsPanelProps {
  visible: boolean;
  player: Player;
  onLogout: () => void;
  onClose?: () => void;
}

export default function SettingsPanel({ visible, player, onLogout, onClose }: SettingsPanelProps) {
  const { showLoadingPopup } = useLoadingContext();

  const handleAdvanceDay = async () => {
    showLoadingPopup('loading', 'Advancing to next day...');

    try {
      const response = await makeAuthenticatedRequest('/api/nextDay', {
        method: 'POST',
      });

      if (response.ok) {
        showLoadingPopup('success', 'Day advanced successfully! Please refresh the page', () => {
          window.location.reload();
        });
      } else {
        const errorData = await response.json();
        showLoadingPopup('error', `Failed to advance day: ${errorData.error || 'Unknown error'}.`);
      }
    } catch {
      showLoadingPopup('error', 'Failed to advance day: Network error.');
    }
  };
  
  return (
    <SidePanel side="right" visible={visible} onClose={onClose}>
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
  );
}