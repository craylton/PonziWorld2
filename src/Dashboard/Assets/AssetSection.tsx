import type { AvailableAsset } from '../../models/AvailableAsset';
import AssetProvider from '../../contexts/AssetContext';
import AssetList from './AssetList';

interface AssetSectionProps {
  availableAssets: AvailableAsset[];
  onRefreshBank: () => Promise<void>;
}

export default function AssetSection({ availableAssets, onRefreshBank }: AssetSectionProps) {
  const loadInvestedAssets = async (): Promise<AvailableAsset[]> => {
    return availableAssets.filter(asset => asset.isInvestedOrPending);
  };

  const loadAvailableAssets = async (): Promise<AvailableAsset[]> => {
    return availableAssets.filter(asset => !asset.isInvestedOrPending);
  };

  return (
    <AssetProvider refreshBank={onRefreshBank}>
      <AssetList
        title="Your Assets"
        onLoad={loadInvestedAssets}
        isExpandedByDefault={true}
      />
      
      <AssetList
        title="Available Assets"
        onLoad={loadAvailableAssets}
        isExpandedByDefault={false}
      />
    </AssetProvider>
  );
}
