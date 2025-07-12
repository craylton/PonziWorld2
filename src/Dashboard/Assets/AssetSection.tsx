import type { AssetType } from '../../models/AssetType';
import type { Asset } from './Asset';
import { makeAuthenticatedRequest } from '../../auth';
import AssetList from './AssetList';

const generateRandomDataPoints = (length = 8): number[] => {
  const dataPoints: number[] = [];
  let currentValue = 100;

  for (let i = 0; i < length; i++) {
    dataPoints.push(currentValue);
    const factor = 0.8 + Math.random() * 0.5; // 0.8 to 1.3
    currentValue = Math.round(currentValue * factor);
  }

  return dataPoints;
};

interface AssetSectionProps {
  bankAssets: Asset[];
}

export default function AssetSection({ bankAssets }: AssetSectionProps) {
  // Convert asset types to assets with 0 amount, filtering out ones we've already invested in
  const getFilteredAssetTypes = (allAssetTypes: AssetType[]): Asset[] => {
    if (!allAssetTypes.length) return [];

    const investedAssetTypes = new Set(bankAssets.map(asset => asset.assetType));
    return allAssetTypes
      .filter(assetType => !investedAssetTypes.has(assetType.name))
      .map(assetType => ({
        assetType: assetType.name,
        assetTypeId: assetType.id,
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

  const getInvestedAssetTypes = async (): Promise<Asset[]> => {
    return bankAssets.map(asset => ({
      assetType: asset.assetType,
      assetTypeId: asset.assetTypeId,
      amount: asset.amount,
      dataPoints: generateRandomDataPoints()
    }))
  };

  return (
    <>
      <AssetList
        title="Your Assets"
        onLoad={getInvestedAssetTypes}
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
