import type { Asset } from './Asset';
import InvestedAssetSummary from './InvestedAssetSummary';
import UninvestedAssetSummary from './UninvestedAssetSummary';

interface ConditionalAssetSummaryProps {
  asset: Asset;
  historicalValues: number[];
  isInvested: boolean;
}

export default function ConditionalAssetSummary({ asset, historicalValues, isInvested }: ConditionalAssetSummaryProps) {
  if (isInvested) {
    return <InvestedAssetSummary asset={asset} historicalValues={historicalValues} />;
  } else {
    return <UninvestedAssetSummary asset={asset} historicalValues={historicalValues} />;
  }
}
