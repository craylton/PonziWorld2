import AssetSummaryBase from './AssetSummaryBase';
import type { Asset } from './Asset';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';

interface AssetSummaryProps {
  asset: Asset;
  historicalValues: number[];
}

export default function InvestedAssetSummary({ asset, historicalValues }: AssetSummaryProps) {
  const hasPendingAmount = asset.pendingAmount !== 0;
  
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{asset.assetType}</div>
          <div className="asset-list__amount">
            {hasPendingAmount ? (
              <>
                {formatCurrency(asset.amount)} {asset.pendingAmount > 0 ? '+' : '-'} {' '}
                <span 
                  className={`asset-list__pending ${asset.pendingAmount > 0 ? 'asset-list__pending--positive' : 'asset-list__pending--negative'}`}
                >
                  {formatCurrency(Math.abs(asset.pendingAmount))}
                </span>
              </>
            ) : (
              formatCurrency(asset.amount)
            )}
          </div>
        </div>
        <AssetSummaryBase asset={asset} historicalValues={historicalValues} />
      </div>
    </>
  );
}
