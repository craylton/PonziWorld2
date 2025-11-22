import type { AvailableAsset } from './AvailableAsset';
import type { InvestmentDetailsResponse } from './AssetDetails';

export interface AssetWithDetails {
    asset: AvailableAsset;
    details: InvestmentDetailsResponse;
}
