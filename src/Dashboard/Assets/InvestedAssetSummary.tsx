import AssetHistoricalPerformanceDisplay from './AssetHistoricalPerformanceDisplay';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';

interface AssetSummaryProps {
  investment: InvestmentDetailsResponse;
}

export default function InvestedAssetSummary({ investment }: AssetSummaryProps) {
  const hasPendingAmount = investment.pendingAmount !== 0;
  
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{investment.targetAssetName}</div>
          <div className="asset-list__amount">
            {hasPendingAmount ? (
              <>
                {formatCurrency(investment.investedAmount)} {investment.pendingAmount > 0 ? '+' : '-'}
                <span 
                  className={`asset-list__pending ${investment.pendingAmount > 0 ? 'asset-list__pending--positive' : 'asset-list__pending--negative'}`}
                >
                  {formatCurrency(Math.abs(investment.pendingAmount))}
                </span>
              </>
            ) : (
              formatCurrency(investment.investedAmount)
            )}
          </div>
        </div>
        <AssetHistoricalPerformanceDisplay investment={investment} />
      </div>
    </>
  );
}
