import type { AssetType } from '../../models/AssetType';
import type { Asset as ModelAsset } from '../../models/Asset';
import type { Asset } from './Asset';
import type { PendingTransaction } from '../../models/PendingTransaction';
import { makeAuthenticatedRequest } from '../../auth';
import { useBankContext } from '../../contexts/useBankContext';
import AssetProvider from '../../contexts/AssetContext';
import AssetList from './AssetList';
import { useCallback } from 'react';

interface AssetSectionProps {
  bankAssets: ModelAsset[];
}

export default function AssetSection({ bankAssets }: AssetSectionProps) {
  const { bankId } = useBankContext();

  // Fetch asset types from backend
  const fetchAssetTypes = useCallback(async (): Promise<AssetType[]> => {
    try {
      const res = await makeAuthenticatedRequest('/api/assetTypes');
      if (res.ok) return res.json();
      console.error('Failed to load asset types');
      return [];
    } catch (error) {
      console.error('Error loading asset types:', error);
      return [];
    }
  }, []);

  // Fetch pending transactions for the bank
  const fetchPendingTransactions = useCallback(async (): Promise<PendingTransaction[]> => {
    try {
      const response = await makeAuthenticatedRequest(`/api/pendingTransactions/${bankId}`);
      if (response.ok)
        return response.json();
      console.error('Failed to load pending transactions');
      return [];
    } catch (error) {
      console.error('Error loading pending transactions:', error);
      return [];
    }
  }, [bankId]);


  // Convert asset types to assets with 0 amount, filtering out ones we've already invested in
  const getFilteredAssetTypes = useCallback((
    allAssetTypes: AssetType[],
    pendingTransactions: PendingTransaction[]
  ): Asset[] => {
    if (!allAssetTypes.length) return [];

    const investedAssetTypes = new Set(bankAssets.map(asset => asset.assetType));
    const pendingAssetTypeIds = new Set(pendingTransactions.map(pt => pt.assetId));
    return allAssetTypes
      .filter(assetType => !investedAssetTypes.has(assetType.name) && !pendingAssetTypeIds.has(assetType.id))
      .map(assetType => ({
        assetType: assetType.name,
        assetTypeId: assetType.id,
        amount: 0,
        pendingAmount: 0
      }));
  }, [bankAssets]);

  const fetchAvailableAssetTypes = useCallback(async (): Promise<Asset[]> => {
    const [assetTypes, pendingTransactions] = await Promise.all([
      fetchAssetTypes(),
      fetchPendingTransactions(),
    ]);

    return getFilteredAssetTypes(assetTypes, pendingTransactions);
  }, [fetchAssetTypes, fetchPendingTransactions, getFilteredAssetTypes]);

  const getInvestedAssetTypes = useCallback(async (): Promise<Asset[]> => {
    const [pendingTransactions, allAssetTypes] = await Promise.all([
      fetchPendingTransactions(),
      fetchAssetTypes(),
    ]);

    const assetTypeMap = new Map(allAssetTypes.map(at => [at.id, at.name]));

    // Start with invested assets and add pending amounts
    const investedAssets = bankAssets.map(asset => {
      const pendingTransaction = pendingTransactions.find(pt => pt.assetId === asset.assetTypeId);
      return {
        assetType: asset.assetType,
        assetTypeId: asset.assetTypeId,
        amount: asset.amount,
        pendingAmount: pendingTransaction?.amount || 0,
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
        pendingAmount: pt.amount,
      }));

    return [...investedAssets, ...pendingOnlyAssets];
  }, [fetchAssetTypes, fetchPendingTransactions, bankAssets]);

  return (
    <AssetProvider>
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
    </AssetProvider>
  );
}
