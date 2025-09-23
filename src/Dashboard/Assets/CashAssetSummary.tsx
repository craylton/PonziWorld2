import { useEffect } from 'react';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';
import { useAssetContext } from '../../contexts/useAssetContext';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';

interface CashAssetSummaryProps {
  investment: InvestmentDetailsResponse;
}

export default function CashAssetSummary({ investment }: CashAssetSummaryProps) {
  const { setCashBalance } = useAssetContext();
  const hasPendingAmount = investment.pendingAmount !== 0;
  
  // Update cash balance in context whenever the cash asset data changes
  useEffect(() => {
    const totalCashBalance = investment.investedAmount + investment.pendingAmount;
    setCashBalance(totalCashBalance);
  }, [investment.investedAmount, investment.pendingAmount, setCashBalance]);
  
  return (
    <>
      <div className="asset-list__item asset-list__item--cash">
        <div className="asset-list__content--cash">
          <div className="asset-list__type">{investment.targetAssetName}:</div>
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
      </div>
    </>
  );
}
