import type { AssetType } from '../../models/AssetType';
import type { Asset } from './Asset';
import type { PendingTransaction } from '../../models/PendingTransaction';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';
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
  const { bankId } = useBankContext();

  // Fetch pending transactions for the bank
  const fetchPendingTransactions = async (): Promise<PendingTransaction[]> => {
    try {
      const response = await makeAuthenticatedRequest(`/api/pendingTransactions/${bankId}`);
      if (response.ok) {
        return await response.json();
      } else {
        console.error('Failed to load pending transactions');
        return [];
      }
    } catch (error) {
      console.error('Error loading pending transactions:', error);
      return [];
    }
  };

  // Convert asset types to assets with 0 amount, filtering out ones we've already invested in
  const getFilteredAssetTypes = (allAssetTypes: AssetType[], pendingTransactions: PendingTransaction[]): Asset[] => {
    if (!allAssetTypes.length) return [];

    const investedAssetTypes = new Set(bankAssets.map(asset => asset.assetType));
    const pendingAssetTypeIds = new Set(pendingTransactions.map(pt => pt.assetId));
    
    return allAssetTypes
      .filter(assetType => !investedAssetTypes.has(assetType.name) && !pendingAssetTypeIds.has(assetType.id))
      .map(assetType => ({
        assetType: assetType.name,
        assetTypeId: assetType.id,
        amount: 0
      }));
  };

  const fetchAvailableAssetTypes = async (): Promise<Asset[]> => {
    try {
      const [assetTypesResponse, pendingTransactions] = await Promise.all([
        makeAuthenticatedRequest('/api/assetTypes'),
        fetchPendingTransactions()
      ]);
      
      if (assetTypesResponse.ok) {
        const assetTypes: AssetType[] = await assetTypesResponse.json();
        // create assets with random dataPoints
        return getFilteredAssetTypes(assetTypes, pendingTransactions).map(asset => ({
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
    try {
      const [pendingTransactions, assetTypesResponse] = await Promise.all([
        fetchPendingTransactions(),
        makeAuthenticatedRequest('/api/assetTypes')
      ]);

      const allAssetTypes: AssetType[] = assetTypesResponse.ok ? await assetTypesResponse.json() : [];
      const assetTypeMap = new Map(allAssetTypes.map(at => [at.id, at.name]));

      // Start with invested assets and add pending amounts
      const investedAssets = bankAssets.map(asset => {
        const pendingTransaction = pendingTransactions.find(pt => pt.assetId === asset.assetTypeId);
        return {
          assetType: asset.assetType,
          assetTypeId: asset.assetTypeId,
          amount: asset.amount,
          dataPoints: generateRandomDataPoints(),
          pendingAmount: pendingTransaction?.amount || 0
        };
      });

      // Add assets that have pending transactions but no current investment
      const investedAssetIds = new Set(bankAssets.map(asset => asset.assetTypeId));
      const pendingOnlyAssets = pendingTransactions
        .filter(pt => !investedAssetIds.has(pt.assetId))
        .map(pt => ({
          assetType: assetTypeMap.get(pt.assetId) || 'Unknown Asset',
          assetTypeId: pt.assetId,
          amount: 0,
          dataPoints: generateRandomDataPoints(),
          pendingAmount: pt.amount
        }));

      return [...investedAssets, ...pendingOnlyAssets];
    } catch (error) {
      console.error('Error loading invested assets:', error);
      // Fallback to just the bank assets without pending transaction info
      return bankAssets.map(asset => ({
        assetType: asset.assetType,
        assetTypeId: asset.assetTypeId,
        amount: asset.amount,
        dataPoints: generateRandomDataPoints()
      }));
    }
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
