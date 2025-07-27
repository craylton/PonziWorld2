import { useEffect } from 'react';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';
import { useAssetContext } from '../../contexts/useAssetContext';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';

interface CashAssetSummaryProps {
  asset: InvestmentDetailsResponse;
}

export default function CashAssetSummary({ asset }: CashAssetSummaryProps) {
  const { setCashBalance } = useAssetContext();
  const hasPendingAmount = asset.pendingAmount !== 0;
  
  // Update cash balance in context whenever the cash asset data changes
  useEffect(() => {
    const totalCashBalance = asset.investedAmount + asset.pendingAmount;
    setCashBalance(totalCashBalance);
  }, [asset.investedAmount, asset.pendingAmount, setCashBalance]);
  
  return (
    <>
      <div className="asset-list__item asset-list__item--cash">
        <div className="asset-list__content">
          <div className="asset-list__type">{asset.name}</div>
          <div className="asset-list__amount">
            {hasPendingAmount ? (
              <>
                {formatCurrency(asset.investedAmount)} {asset.pendingAmount > 0 ? '+' : '-'}
                <span 
                  className={`asset-list__pending ${asset.pendingAmount > 0 ? 'asset-list__pending--positive' : 'asset-list__pending--negative'}`}
                >
                  {formatCurrency(Math.abs(asset.pendingAmount))}
                </span>
              </>
            ) : (
              formatCurrency(asset.investedAmount)
            )}
          </div>
        </div>
      </div>
    </>
  );
}
