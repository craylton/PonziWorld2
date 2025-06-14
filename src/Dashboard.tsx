import React, { useState } from 'react';
import DashboardHeader from './DashboardHeader';
import './Dashboard.css';

interface DashboardProps {
  username: string;
}

// TODO: Replace with real data from backend
const DUMMY_BANK_NAME = 'Ponzi National Bank';
const DUMMY_CLAIMED_CAPITAL = 1000000;
const DUMMY_ACTUAL_CAPITAL = 250000;

export default function Dashboard({ username }: DashboardProps) {
    const [showLeftPanel, setShowLeftPanel] = useState(false);
    const [showRightPanel, setShowRightPanel] = useState(false);
    const mainContent = `Welcome to the dashboard, ${username}!`;
  return (
    <div className="dashboard-root">
      <DashboardHeader
        bankName={DUMMY_BANK_NAME}
        claimedCapital={DUMMY_CLAIMED_CAPITAL}
        actualCapital={DUMMY_ACTUAL_CAPITAL}
      />
      <div className="dashboard-layout">
        {/* Sliding left panel with hamburger/X button inside */}
        <aside
          className={`dashboard-panel dashboard-panel--left${showLeftPanel ? ' dashboard-panel--visible' : ''}`}
        >
          <button
            className={`dashboard-hamburger dashboard-hamburger--panel${showLeftPanel ? ' dashboard-hamburger--open' : ''}`}
            aria-label={showLeftPanel ? 'Close left panel' : 'Open left panel'}
            onClick={() => setShowLeftPanel((v) => !v)}
          >
            <span />
            <span />
            <span />
          </button>
          Left Panel
        </aside>
        <main className="dashboard-main">
          {/* Hamburger for left panel (mobile only, floats over main) */}
          <button
            className={`dashboard-hamburger dashboard-hamburger--float dashboard-hamburger--left${showLeftPanel ? ' dashboard-hamburger--hidden' : ''}`}
            aria-label="Open left panel"
            onClick={() => setShowLeftPanel(true)}
          >
            <span />
            <span />
            <span />
          </button>
          {mainContent}
          {/* Hamburger for right panel (mobile only, floats over main) */}
          <button
            className={`dashboard-hamburger dashboard-hamburger--float dashboard-hamburger--right${showRightPanel ? ' dashboard-hamburger--hidden' : ''}`}
            aria-label="Open right panel"
            onClick={() => setShowRightPanel(true)}
          >
            <span />
            <span />
            <span />
          </button>
        </main>
        {/* Sliding right panel with hamburger/X button inside */}
        <aside
          className={`dashboard-panel dashboard-panel--right${showRightPanel ? ' dashboard-panel--visible' : ''}`}
        >
          <button
            className={`dashboard-hamburger dashboard-hamburger--panel${showRightPanel ? ' dashboard-hamburger--open' : ''}`}
            aria-label={showRightPanel ? 'Close right panel' : 'Open right panel'}
            onClick={() => setShowRightPanel((v) => !v)}
          >
            <span />
            <span />
            <span />
          </button>
          Right Panel
        </aside>
      </div>
    </div>
  );
}
