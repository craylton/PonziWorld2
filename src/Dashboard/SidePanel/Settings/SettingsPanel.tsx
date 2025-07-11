import SidePanel from '../SidePanel';
import type { Player } from '../../../models/User';
import { makeAuthenticatedRequest } from '../../../auth';

interface SettingsPanelProps {
  visible: boolean;
  player: Player;
  onLogout: () => void;
  onClose?: () => void;
}

const handleAdvanceDay = async () => {
  try {
    const response = await makeAuthenticatedRequest('/api/nextDay', {
      method: 'POST',
    });

    if (response.ok) {
      // Do nothing, maybe show a success message one day
    } else {
      const errorData = await response.json();
      alert(`Failed to advance day: ${errorData.error || 'Unknown error'}`);
    }
  } catch {
    alert('Failed to advance day: Network error');
  }
};

export default function SettingsPanel({ visible, player, onLogout, onClose }: SettingsPanelProps) {
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
