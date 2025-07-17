import AssetSummaryBase from './AssetSummaryBase';
import type { Asset } from './Asset';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';
import { useState } from 'react';

interface AssetSummaryProps {
  asset: Asset;
}

export default function InvestedAssetSummary({ asset }: AssetSummaryProps) {
  const [investedAmount, setInvestedAmount] = useState<number>(asset.amount);
  const [pendingAmount, setPendingAmount] = useState<number>(asset.pendingAmount);
  const [isLoading, setIsLoading] = useState(true);

  const handleAssetDetailsLoaded = (loadedInvestedAmount: number, loadedPendingAmount: number) => {
    setInvestedAmount(loadedInvestedAmount);
    setPendingAmount(loadedPendingAmount);
    setIsLoading(false);
  };

  const hasPendingAmount = pendingAmount !== 0;
  
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{asset.assetType}</div>
          <div className="asset-list__amount">
            {isLoading ? (
              <span style={{ color: 'rgba(255, 255, 255, 0.6)' }}>Loading...</span>
            ) : (
              <>
                {hasPendingAmount ? (
                  <>
                    {formatCurrency(investedAmount)} {pendingAmount > 0 ? '+' : '-'} {' '}
                    <span 
                      className={`asset-list__pending ${pendingAmount > 0 ? 'asset-list__pending--positive' : 'asset-list__pending--negative'}`}
                    >
                      {formatCurrency(Math.abs(pendingAmount))}
                    </span>
                  </>
                ) : (
                  formatCurrency(investedAmount)
                )}
              </>
            )}
          </div>
        </div>
        <AssetSummaryBase asset={asset} onAssetDetailsLoaded={handleAssetDetailsLoaded} />
      </div>
    </>
  );
}
