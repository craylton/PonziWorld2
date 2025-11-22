import './BottomPanel.css';

interface BottomTabsProps {
  activeTab: 'today' | 'tomorrow' | null;
  onTabClick: (tab: 'today' | 'tomorrow') => void;
}

export default function BottomTabs({ activeTab, onTabClick }: BottomTabsProps) {
  return (
    <div className="dashboard-bottom-tabs">
      <button
        type="button"
        className={`dashboard-bottom-tab${activeTab === 'today' ? ' dashboard-bottom-tab--active' : ''}`}
        onClick={() => onTabClick('today')}
      >
        Today
      </button>
      <button
        type="button"
        className={`dashboard-bottom-tab${activeTab === 'tomorrow' ? ' dashboard-bottom-tab--active' : ''}`}
        onClick={() => onTabClick('tomorrow')}
      >
        Tomorrow
      </button>
    </div>
  );
}
