import type { AvailableAsset } from '../../models/AvailableAsset';
import type { AssetDetailsResponse } from '../../models/AssetDetails';
import InvestedAssetSummary from './InvestedAssetSummary';
import UninvestedAssetSummary from './UninvestedAssetSummary';
import { useState, useEffect } from 'react';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';

interface AssetSummaryProps {
    availableAsset: AvailableAsset;
}

export default function AssetSummary({ availableAsset }: AssetSummaryProps) {
    const { bankId } = useBankContext();
    const [assetDetails, setAssetDetails] = useState<AssetDetailsResponse | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchAssetDetails = async () => {
            setIsLoading(true);
            try {
                const response = await makeAuthenticatedRequest(
                    `/api/asset/${availableAsset.assetTypeId}/${bankId}`
                );
                if (response.ok) {
                    const data: AssetDetailsResponse = await response.json();
                    setAssetDetails(data);
                } else {
                    console.error('Failed to fetch asset details for asset:', availableAsset.assetType);
                }
            } catch (error) {
                console.error('Error fetching asset details:', error);
            } finally {
                setIsLoading(false);
            }
        };

        fetchAssetDetails();
    }, [availableAsset.assetTypeId, bankId, availableAsset.assetType]);

    if (isLoading || !assetDetails) {
        return (
            <div className="asset-list__item">
                <div className="asset-list__content">
                    <div className="asset-list__type">{availableAsset.assetType}</div>
                    <div className="asset-list__amount">Loading...</div>
                </div>
            </div>
        );
    }

    return availableAsset.isInvestedOrPending ? (
        <InvestedAssetSummary asset={assetDetails!} />
    ) : (
        <UninvestedAssetSummary asset={assetDetails!} />
    );
}
