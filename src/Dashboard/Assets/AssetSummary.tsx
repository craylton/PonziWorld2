import type { AvailableAsset } from '../../models/AvailableAsset';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';
import InvestedAssetSummary from './InvestedAssetSummary';
import UninvestedAssetSummary from './UninvestedAssetSummary';
import CashAssetSummary from './CashAssetSummary';
import { useState, useEffect, useCallback } from 'react';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';
import { useAssetContext } from '../../contexts/useAssetContext';

interface AssetSummaryProps {
    availableAsset: AvailableAsset;
}

export default function AssetSummary({ availableAsset }: AssetSummaryProps) {
    const { bankId } = useBankContext();
    const { registerRefreshCallback, unregisterRefreshCallback } = useAssetContext();
    const [assetDetails, setAssetDetails] = useState<InvestmentDetailsResponse | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    const fetchAssetDetails = useCallback(async () => {
        setIsLoading(true);
        try {
            const response = await makeAuthenticatedRequest(
                `/api/investment/${availableAsset.assetId}/${bankId}`
            );
            if (response.ok) {
                const data: InvestmentDetailsResponse = await response.json();
                setAssetDetails(data);
            } else {
                console.error('Failed to fetch asset details for asset:', availableAsset.assetName);
            }
        } catch (error) {
            console.error('Error fetching asset details:', error);
        } finally {
            setIsLoading(false);
        }
    }, [availableAsset.assetId, bankId, availableAsset.assetName]);

    useEffect(() => {
        fetchAssetDetails();
    }, [fetchAssetDetails]);

    // Register for refresh callbacks
    useEffect(() => {
        registerRefreshCallback(fetchAssetDetails);
        return () => {
            unregisterRefreshCallback(fetchAssetDetails);
        };
    }, [registerRefreshCallback, unregisterRefreshCallback, fetchAssetDetails]);

    if (isLoading || !assetDetails) {
        return (
            <div className="asset-list__item">
                <div className="asset-list__content">
                    <div className="asset-list__type">{availableAsset.assetName}</div>
                    <div className="asset-list__amount">Loading...</div>
                </div>
            </div>
        );
    }

    if (assetDetails.name === 'Cash') {
        return <CashAssetSummary asset={assetDetails} />;
    }

    return availableAsset.isInvestedOrPending ? (
        <InvestedAssetSummary asset={assetDetails!} />
    ) : (
        <UninvestedAssetSummary asset={assetDetails!} />
    );
}
