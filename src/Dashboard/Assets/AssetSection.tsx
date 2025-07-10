import type { AssetType } from '../../models/AssetType';
import type { Asset } from './Asset';
import { makeAuthenticatedRequest } from '../../auth';
import AssetList from './AssetList';

const generateRandomDataPoints = (length = 8): number[] => {
  const dataPoints: number[] = [];
  let currentValue = 1000;
  
  for (let i = 0; i < length; i++) {
    dataPoints.push(currentValue);
    const factor = 0.9 + Math.random() * 0.3; // 0.9 to 1.2
    currentValue = Math.round(currentValue * factor);
  }
  
  return dataPoints;
};

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
