import type { AvailableAsset } from '../../models/AvailableAsset';
import AssetProvider from '../../contexts/AssetContext';
import AssetList from './AssetList';
import { useEffect, useState } from 'react';

interface AssetSectionProps {
  availableAssets: AvailableAsset[];
  onRefreshBank: () => Promise<void>;
}

export default function AssetSection({ availableAssets, onRefreshBank }: AssetSectionProps) {
  const [investedAssets, setInvestedAssets] = useState<AvailableAsset[]>([]);
  const [uninvestedAssets, setUninvestedAssets] = useState<AvailableAsset[]>([]);

  useEffect(() => {
    const invested = availableAssets.filter(asset => asset.isInvestedOrPending);
    const uninvested = availableAssets.filter(asset => !asset.isInvestedOrPending);
    setInvestedAssets(invested);
    setUninvestedAssets(uninvested);
  }, [availableAssets]);

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
