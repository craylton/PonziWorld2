import { useState, useMemo } from 'react';
import './InvestorList.css';
import type { Investor } from './Investor';
import InvestorSummary from './InvestorSummary';

type SortOption = 'alphabetical' | 'investment';

interface InvestorListProps {
  investors: Investor[];
}

export default function InvestorList({ investors }: InvestorListProps) {
  const [sortBy, setSortBy] = useState<SortOption>('investment');

  const sortedInvestors = useMemo(() => {
    const sorted = [...investors];
    
    if (sortBy === 'alphabetical') {
      sorted.sort((a, b) => a.name.localeCompare(b.name));
    } else {
      sorted.sort((a, b) => b.amountInvested - a.amountInvested);
    }
    
    return sorted;
  }, [investors, sortBy]);

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
        {sortedInvestors.length === 0 ? (
          <div className="investor-list__empty">No investors yet</div>
        ) : (
          sortedInvestors.map((investor) => (
            <InvestorSummary key={investor.id} investor={investor} />
          ))
        )}
      </div>
    </div>
  );
}
