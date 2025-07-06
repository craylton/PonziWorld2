import SidePanelButton from './SidePanel/SidePanelButton';

interface InvestorsButtonProps {
  isLeftPanelOpen: boolean;
  onClick: () => void;
}

export default function InvestorsButton({ isLeftPanelOpen, onClick }: InvestorsButtonProps) {
  return (
    <SidePanelButton
      iconType="hamburger"
      shouldAllowClose={isLeftPanelOpen}
      onClick={onClick}
      ariaLabel="Open left panel"
      className={`dashboard-sidepanel-button--left`}
    />
  );
}
