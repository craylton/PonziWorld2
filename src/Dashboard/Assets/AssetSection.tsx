import type { AssetType } from '../../models/AssetType';
import type { Asset } from './Asset';
import { makeAuthenticatedRequest } from '../../auth';
import AssetList from './AssetList';

// Helper to generate random data points for visualization
const generateRandomDataPoints = (length = 7): number[] =>
  Array.from({ length }, () => Math.floor(Math.random() * 20) + 1);

interface AssetSectionProps {
  bankAssets: Asset[];
}

export default function AssetSection({ bankAssets }: AssetSectionProps) {
  // Convert asset types to assets with 0 amount, filtering out existing ones
  const getFilteredAssetTypes = (assetTypes: AssetType[]): Asset[] => {
    if (!assetTypes.length) return [];

    const existingAssetTypes = new Set(bankAssets.map(asset => asset.assetType));
    return assetTypes
      .filter(assetType => !existingAssetTypes.has(assetType.name))
      .map(assetType => ({
        assetType: assetType.name,
        amount: 0
      }));
  };

  const fetchAvailableAssetTypes = async (): Promise<Asset[]> => {
    try {
      const response = await makeAuthenticatedRequest('/api/assetTypes');
      if (response.ok) {
        const assetTypes: AssetType[] = await response.json();
        // create assets with random dataPoints
        return getFilteredAssetTypes(assetTypes).map(asset => ({
          ...asset,
          dataPoints: generateRandomDataPoints()
        }));
      } else {
        console.error('Failed to load asset types');
        return [];
      }
    } catch (error) {
      console.error('Error loading asset types:', error);
      return [];
    }
  };

  return (
    <>
      <AssetList
        title="Your Assets"
        onLoad={() =>
          Promise.resolve(
            bankAssets.map(asset => ({
              assetType: asset.assetType,
              amount: asset.amount,
              dataPoints: generateRandomDataPoints()
            }))
          )
        }
        isExpandedByDefault
      />

      <AssetList
        title="Available Assets"
        onLoad={fetchAvailableAssetTypes}
        isExpandedByDefault={false}
      />
    </>
  );
}
