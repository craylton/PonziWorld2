import { useState, useEffect } from 'react';
import SidePanel from './SidePanel';
import './Dashboard.css';
import DashboardHeader from './DashboardHeader';
import HamburgerButton from './HamburgerButton';
import InvestorList from './InvestorList';
import CogButton from './CogButton';

interface DashboardProps {
  username: string;
}

// TODO: Replace with real data from backend
const DUMMY_BANK_NAME = 'Ponzi National Bank';
const DUMMY_CLAIMED_CAPITAL = 1000000;
const DUMMY_ACTUAL_CAPITAL = 250000;

// Hook to get window width
function useWindowWidth() {
  const [width, setWidth] = useState(window.innerWidth);
  useEffect(() => {
    const handleResize = () => setWidth(window.innerWidth);
    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
  }, []);
  return width;
}

export default function Dashboard({ username }: DashboardProps) {
  const [showLeftPanel, setShowLeftPanel] = useState(false);
  const [showRightPanel, setShowRightPanel] = useState(false);
  const windowWidth = useWindowWidth();
  const mainContent = `Welcome to the dashboard, ${username}!`;

  return (
    <div className="dashboard-root">
      <DashboardHeader
        bankName={DUMMY_BANK_NAME}
        claimedCapital={DUMMY_CLAIMED_CAPITAL}
        actualCapital={DUMMY_ACTUAL_CAPITAL}
      />
      <div className="dashboard-layout">
        {/* Left Side Panel */}
        <SidePanel side="left" visible={showLeftPanel}>
          {/* Only show X button if under 900px and left panel is open */}
          {windowWidth < 900 && showLeftPanel && (
            <HamburgerButton
              shouldAllowClose={true}
              onClick={() => setShowLeftPanel(false)}
              ariaLabel="Close left panel"
            />
          )}
          <InvestorList />
        </SidePanel>
        <main className="dashboard-main">
          {/* Hamburger for left panel (mobile only, floats over main) */}
          {windowWidth < 900 && !showLeftPanel && (
            <HamburgerButton
              shouldAllowClose={false}
              onClick={() => setShowLeftPanel(true)}
              ariaLabel="Open left panel"
              className={`dashboard-hamburger--float dashboard-hamburger--left`}
            />
          )}
          {mainContent}
          {/* Cog for right panel (floats, always visible unless panel is open) */}
          {!showRightPanel && (
            <CogButton
              shouldAllowClose={false}
              onClick={() => setShowRightPanel(true)}
              ariaLabel="Open settings panel"
              className={`dashboard-hamburger--float dashboard-hamburger--right`}
            />
          )}
        </main>
        {/* Right Side Panel (slides in, always overlays on mobile/desktop) */}
        <SidePanel side="right" visible={showRightPanel}>
          {/* Always show X button inside right panel when open */}
          <CogButton
            shouldAllowClose={true}
            onClick={() => setShowRightPanel(false)}
            ariaLabel="Close settings panel"
            className="dashboard-hamburger--panel-close"
          />
          {/* TODO: Add settings content here */}
        </SidePanel>
      </div>
    </div>
  );
}
