import type { AvailableAsset } from '../../models/AvailableAsset';
import AssetProvider from '../../contexts/AssetContext';
import AssetList from './AssetList';

interface AssetSectionProps {
  availableAssets: AvailableAsset[];
}

export default function AssetSection({ availableAssets }: AssetSectionProps) {
  const loadInvestedAssets = async (): Promise<AvailableAsset[]> => {
    return availableAssets.filter(asset => asset.isInvestedOrPending);
  };

  const loadAvailableAssets = async (): Promise<AvailableAsset[]> => {
    return availableAssets.filter(asset => !asset.isInvestedOrPending);
  };

  return (
    <AssetProvider>
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
