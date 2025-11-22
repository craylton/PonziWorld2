import type { AvailableAsset } from '../../models/AvailableAsset';
import type { InvestmentDetailsResponse } from '../../models/AssetDetails';
import type { AssetWithDetails } from '../../models/AssetWithDetails';
import AssetProvider from '../../contexts/AssetContext';
import AssetList from './AssetList';
import { useEffect, useState } from 'react';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';

interface AssetSectionProps {
  availableAssets: AvailableAsset[];
  onRefreshBank: () => Promise<void>;
}

async function fetchAssetDetails(asset: AvailableAsset, bankId: string): Promise<AssetWithDetails | null> {
  const response = await makeAuthenticatedRequest(`/api/investment/${asset.assetId}/${bankId}`);
  if (response.ok) {
    const details: InvestmentDetailsResponse = await response.json();
    return { asset, details };
  }
  return null;
}

function sortInvestedAssets(assets: AssetWithDetails[]): AssetWithDetails[] {
  return assets.sort((a, b) => {
    if (a.details.targetAssetName === 'Cash') return -1;
    if (b.details.targetAssetName === 'Cash') return 1;
    return b.details.investedAmount - a.details.investedAmount;
  });
}

export default function AssetSection({ availableAssets, onRefreshBank }: AssetSectionProps) {
  const { bankId } = useBankContext();
  const [investedAssets, setInvestedAssets] = useState<AssetWithDetails[]>([]);
  const [uninvestedAssets, setUninvestedAssets] = useState<AssetWithDetails[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    async function loadAssets() {
      setIsLoading(true);
      try {
        const fetchPromises = availableAssets.map(asset => fetchAssetDetails(asset, bankId));
        const results = await Promise.all(fetchPromises);
        const validResults = results.filter((result): result is AssetWithDetails => result !== null);

        const invested = validResults.filter(item => item.asset.isInvestedOrPending);
        const uninvested = validResults.filter(item => !item.asset.isInvestedOrPending);

        setInvestedAssets(sortInvestedAssets(invested));
        setUninvestedAssets(uninvested);
      } catch (error) {
        console.error('Error fetching asset details:', error);
      } finally {
        setIsLoading(false);
      }
    }

    loadAssets();
  }, [availableAssets, bankId]);

  if (isLoading) {
    return (
      <AssetProvider refreshBank={onRefreshBank}>
        <div className="asset-list">
          <div className="asset-list__items">
            <div className="asset-list__item">
              <div className="asset-list__content">
                <div className="asset-list__amount">Loading assets...</div>
              </div>
            </div>
          </div>
        </div>
      </AssetProvider>
    );
  }

  return (
    <AssetProvider refreshBank={onRefreshBank}>
      <AssetList
        title="Your Investments"
        isExpandedByDefault={true}
        assets={investedAssets}
      />

      <AssetList
        title="Available Assets"
        assets={uninvestedAssets}
        isExpandedByDefault={false}
      />
    </AssetProvider>
  );
}
