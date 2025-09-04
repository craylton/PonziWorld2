import { useEffect } from 'react';
import './AssetList.css';
import { formatCurrencyFromString } from '../../utils/currency';
import { addMoney, parseMoney } from '../../utils/money';
import { useAssetContext } from '../../contexts/useAssetContext';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';

interface CashAssetSummaryProps {
  investment: InvestmentDetailsResponse;
}

export default function CashAssetSummary({ investment }: CashAssetSummaryProps) {
  const { setCashBalance } = useAssetContext();
  const investedAmount = parseMoney(investment.investedAmount);
  const pendingAmount = parseMoney(investment.pendingAmount);
  const hasPendingAmount = !pendingAmount.isZero();
  
  // Update cash balance in context whenever the cash asset data changes
  useEffect(() => {
    const totalCashBalance = addMoney(investedAmount, pendingAmount);
    setCashBalance(totalCashBalance.toString());
  }, [investment.investedAmount, investment.pendingAmount, setCashBalance, investedAmount, pendingAmount]);
  
  return (
    <>
      <div className="asset-list__item asset-list__item--cash">
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
      </div>
    </>
  );
}
