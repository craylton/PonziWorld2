import AssetSummaryBase from './AssetSummaryBase';
import type { Asset } from './Asset';
import './AssetList.css';

interface AssetSummaryProps {
  asset: Asset;
  historicalValues: number[];
}

export default function UninvestedAssetSummary({ asset, historicalValues }: AssetSummaryProps) {
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{asset.assetType}</div>
        </div>
        <AssetSummaryBase asset={asset} historicalValues={historicalValues} />
      </div>
    </>
  );
}
