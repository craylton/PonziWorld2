import AssetHistoricalPerformanceDisplay from './AssetHistoricalPerformanceDisplay';
import './AssetList.css';
import { formatCurrencyFromString } from '../../utils/currency';
import { parseMoney } from '../../utils/money';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';

interface AssetSummaryProps {
  investment: InvestmentDetailsResponse;
}

export default function InvestedAssetSummary({ investment }: AssetSummaryProps) {
  const pendingAmount = parseMoney(investment.pendingAmount);
  const hasPendingAmount = !pendingAmount.isZero();
  
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{investment.targetAssetName}</div>
          <div className="asset-list__amount">
            {hasPendingAmount ? (
              <>
                {formatCurrencyFromString(investment.investedAmount)} {pendingAmount.isPositive() ? '+' : '-'}
                <span 
                  className={`asset-list__pending ${pendingAmount.isPositive() ? 'asset-list__pending--positive' : 'asset-list__pending--negative'}`}
                >
                  {formatCurrencyFromString(pendingAmount.abs().toString())}
                </span>
              </>
            ) : (
              formatCurrencyFromString(investment.investedAmount)
            )}
          </div>
        </div>
        <AssetHistoricalPerformanceDisplay investment={investment} />
      </div>
    </>
  );
}
