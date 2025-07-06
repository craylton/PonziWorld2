import SidePanelButton from '../SidePanelButton';

interface SettingsButtonProps {
  isRightPanelOpen: boolean;
  onClick: () => void;
}

export default function SettingsButton({ isRightPanelOpen, onClick }: SettingsButtonProps) {
  return (
    <SidePanelButton
      iconType="cog"
      shouldAllowClose={isRightPanelOpen}
      onClick={onClick}
      ariaLabel="Open settings panel"
      className={`dashboard-sidepanel-button--right`}
    />
  );
}
