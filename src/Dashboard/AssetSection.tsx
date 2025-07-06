import { useState } from 'react';
import AssetList from './AssetList/AssetList';
import type { AssetType } from '../models/AssetType';
import type { Asset } from './AssetList/Asset';

interface AssetSectionProps {
  bankAssets: Asset[];
}

export default function AssetSection({ bankAssets }: AssetSectionProps) {
  const [allAssetTypes, setAllAssetTypes] = useState<AssetType[]>([]);
  const [showAllAssets, setShowAllAssets] = useState(false);
  const [isLoadingAssetTypes, setIsLoadingAssetTypes] = useState(false);

  const handleLoadMoreAssets = async () => {
    if (isLoadingAssetTypes) return;
    
    setIsLoadingAssetTypes(true);
    try {
      const response = await fetch('/api/assetTypes');
      if (response.ok) {
        const assetTypes: AssetType[] = await response.json();
        setAllAssetTypes(assetTypes);
        setShowAllAssets(true);
      } else {
        console.error('Failed to load asset types');
      }
    } catch (error) {
      console.error('Error loading asset types:', error);
    } finally {
      setIsLoadingAssetTypes(false);
    }
  };

  // Convert asset types to assets with 0 amount, filtering out existing ones
  const getFilteredAssetTypes = (): Asset[] => {
    if (!allAssetTypes.length) return [];
    
    const existingAssetTypes = new Set(bankAssets.map(asset => asset.assetType));
    return allAssetTypes
      .filter(assetType => !existingAssetTypes.has(assetType.name))
      .map(assetType => ({
        assetType: assetType.name,
        amount: 0
      }));
  };

  return (
    <>
      <AssetList assets={bankAssets} />
      
      {/* More assets button and additional asset list */}
      {!showAllAssets && (
        <div className="dashboard-more-assets-container">
          <button
            onClick={handleLoadMoreAssets}
            disabled={isLoadingAssetTypes}
            className="dashboard-more-assets-button"
          >
            {isLoadingAssetTypes ? 'Loading...' : 'More assets...'}
          </button>
        </div>
      )}
      
      {showAllAssets && (
        <AssetList 
          assets={getFilteredAssetTypes()} 
          showBorder={true}
          title="Available Assets"
        />
      )}
    </>
  );
}
