import AssetSummaryBase from './AssetSummaryBase';
import type { Asset } from './Asset';
import './AssetList.css';

interface AssetSummaryProps {
  asset: Asset;
  historicalValues: number[];
  bankId: string;
}

export default function UninvestedAssetSummary({ asset, historicalValues, bankId }: AssetSummaryProps) {
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{asset.assetType}</div>
        </div>
        <AssetSummaryBase asset={asset} historicalValues={historicalValues} bankId={bankId} />
      </div>
    </>
  );
}
