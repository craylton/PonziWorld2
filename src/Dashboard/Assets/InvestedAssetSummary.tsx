import AssetSummaryBase from './AssetSummaryBase';
import type { Asset } from './Asset';
import './AssetList.css';
import { formatCurrency } from '../../utils/currency';

interface AssetSummaryProps {
  asset: Asset;
  historicalValues: number[];
  bankId: string;
}

export default function InvestedAssetSummary({ asset, historicalValues, bankId }: AssetSummaryProps) {
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{asset.assetType}</div>
          <div className="asset-list__amount">{
            formatCurrency(asset.amount)
          }</div>
        </div>
        <AssetSummaryBase asset={asset} historicalValues={historicalValues} bankId={bankId} />
      </div>
    </>
  );
}
