import AssetSummaryBase from './AssetSummaryBase';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';
import type { AssetDetailsResponse } from '../../models/AssetDetails';

interface AssetSummaryProps {
  asset: AssetDetailsResponse;
}

export default function InvestedAssetSummary({ asset }: AssetSummaryProps) {
  const hasPendingAmount = asset.pendingAmount !== 0;
  
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{asset.name}</div>
          <div className="asset-list__amount">
            {hasPendingAmount ? (
              <>
                {formatCurrency(asset.investedAmount)} {asset.pendingAmount > 0 ? '+' : '-'} {' '}
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
        <AssetSummaryBase asset={asset} />
      </div>
    </>
  );
}
