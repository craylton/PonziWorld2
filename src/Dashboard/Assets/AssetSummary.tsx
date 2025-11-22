import type { AvailableAsset } from '../../models/AvailableAsset';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';
import InvestedAssetSummary from './InvestedAssetSummary';
import UninvestedAssetSummary from './UninvestedAssetSummary';
import CashAssetSummary from './CashAssetSummary';

interface AssetSummaryProps {
    availableAsset: AvailableAsset;
    investmentDetails: InvestmentDetailsResponse;
}

export default function AssetSummary({ availableAsset, investmentDetails }: AssetSummaryProps) {
    if (investmentDetails.targetAssetName === 'Cash') {
        return <CashAssetSummary investment={investmentDetails} />;
    }

    return availableAsset.isInvestedOrPending ? (
        <InvestedAssetSummary investment={investmentDetails} />
    ) : (
        <UninvestedAssetSummary investment={investmentDetails} />
    );
}
