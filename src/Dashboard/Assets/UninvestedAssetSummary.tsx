import AssetHistoricalPerformanceDisplay from './AssetHistoricalPerformanceDisplay';
import './AssetList.css';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';

interface AssetSummaryProps {
  investment: InvestmentDetailsResponse;
}

export default function UninvestedAssetSummary({ investment }: AssetSummaryProps) {
  return (
    <>
      <div className="asset-list__item">
        <div className="asset-list__content">
          <div className="asset-list__type">{investment.targetAssetName}</div>
        </div>
        <AssetHistoricalPerformanceDisplay investment={investment} />
      </div>
    </>
  );
}
