import AssetSummaryBase from './AssetSummaryBase';
import './AssetList.css';
import type { AssetDetailsResponse } from '../../models/AssetDetails';

interface AssetSummaryProps {
  asset: AssetDetailsResponse;
}

export default function UninvestedAssetSummary({ asset }: AssetSummaryProps) {
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{asset.name}</div>
        </div>
        <AssetSummaryBase asset={asset} />
      </div>
    </>
  );
}
