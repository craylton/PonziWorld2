import { useState, useMemo } from 'react';
import './InvestorList.css';
import type { Investor } from './Investor';
import InvestorSummary from './InvestorSummary';

// Dummy data for investors
const DUMMY_INVESTORS: Investor[] = [
  { id: '1', name: 'Alice Johnson', amountInvested: 50000 },
  { id: '2', name: 'Bob Smith', amountInvested: 75000 },
  { id: '3', name: 'Charlie Brown', amountInvested: 25000 },
  { id: '4', name: 'Diana Wells', amountInvested: 120000 },
  { id: '5', name: 'Eric Thompson', amountInvested: 35000 },
  { id: '6', name: 'Fiona Davis', amountInvested: 90000 },
  { id: '7', name: 'George Wilson', amountInvested: 15000 },
  { id: '8', name: 'Helen Martinez', amountInvested: 65000 },
];

type SortOption = 'alphabetical' | 'investment';

export default function InvestorList() {
  const [sortBy, setSortBy] = useState<SortOption>('investment');

  const sortedInvestors = useMemo(() => {
    const sorted = [...DUMMY_INVESTORS];
    
    if (sortBy === 'alphabetical') {
      sorted.sort((a, b) => a.name.localeCompare(b.name));
    } else {
      sorted.sort((a, b) => b.amountInvested - a.amountInvested);
    }
    
    return sorted;
  }, [sortBy]);

  return (
    <div className={`investor-list`}>
      <div className="investor-list__header">
        <h3 className="investor-list__title">Investors</h3>
        <select
          className="investor-list__sort-dropdown"
          value={sortBy}
          onChange={(e) => setSortBy(e.target.value as SortOption)}
        >
          <option value="investment">By Investment</option>
          <option value="alphabetical">Alphabetical</option>
        </select>
      </div>
      
      <div className="investor-list__items">
        {sortedInvestors.map((investor) => (
          <InvestorSummary investor={investor} />
        ))}
      </div>
    </div>
  );
}
