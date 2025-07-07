import { useState } from 'react';
import AssetSummary from './AssetSummary';
import ChevronIcon from '../ChevronIcon';
import type { AssetType } from '../../models/AssetType';
import type { Asset } from './Asset';
import { makeAuthenticatedRequest } from '../../auth';

interface AssetSectionProps {
  bankAssets: Asset[];
}

export default function AssetSection({ bankAssets }: AssetSectionProps) {
  const [allAssetTypes, setAllAssetTypes] = useState<AssetType[]>([]);
  const [showYourAssets, setShowYourAssets] = useState(true); // Expanded by default
  const [showAvailableAssets, setShowAvailableAssets] = useState(false); // Collapsed by default
  const [isLoadingAssetTypes, setIsLoadingAssetTypes] = useState(false);
  const [hasLoadedAssetTypes, setHasLoadedAssetTypes] = useState(false);

  const handleToggleYourAssets = () => {
    setShowYourAssets(!showYourAssets);
  };

  const handleToggleAvailableAssets = async () => {
    if (showAvailableAssets) {
      // Just toggle off if already showing
      setShowAvailableAssets(false);
      return;
    }

    // If we haven't loaded asset types yet, load them
    if (!hasLoadedAssetTypes) {
      setIsLoadingAssetTypes(true);
      try {
        const response = await makeAuthenticatedRequest('/api/assetTypes');
        if (response.ok) {
          const assetTypes: AssetType[] = await response.json();
          setAllAssetTypes(assetTypes);
          setHasLoadedAssetTypes(true);
          setShowAvailableAssets(true);
        } else {
          console.error('Failed to load asset types');
        }
      } catch (error) {
        console.error('Error loading asset types:', error);
      } finally {
        setIsLoadingAssetTypes(false);
      }
    } else {
      // Data already loaded, just show it
      setShowAvailableAssets(true);
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
      {/* Your Assets Section */}
      <div className="asset-list asset-list--bordered">
        <div className="asset-list__header">
          <button
            onClick={handleToggleYourAssets}
            className="dashboard-assets-toggle-button"
          >
            <span className="dashboard-assets-toggle-text">Your Assets</span>
            <ChevronIcon
              className={`dashboard-assets-toggle-icon ${showYourAssets ? 'rotated' : ''}`}
              width={16}
              height={16}
              stroke="currentColor"
              strokeWidth={2}
            />
          </button>
        </div>
        
        {showYourAssets && (
          <div className="asset-list__items">
            {bankAssets.length === 0 ? (
              <div className="asset-list__empty-message">
                You have no assets
              </div>
            ) : (
              bankAssets.map((asset, index) => (
                <AssetSummary key={`${asset.assetType}-${index}`} asset={asset} />
              ))
            )}
          </div>
        )}
      </div>

      {/* Available Assets Section */}
      <div className="asset-list asset-list--bordered">
        <div className="asset-list__header">
          <button
            onClick={handleToggleAvailableAssets}
            disabled={isLoadingAssetTypes}
            className="dashboard-assets-toggle-button"
          >
            <span className="dashboard-assets-toggle-text">
              {isLoadingAssetTypes ? 'Loading...' : 'Available Assets'}
            </span>
            {!isLoadingAssetTypes && (
              <ChevronIcon
                className={`dashboard-assets-toggle-icon ${showAvailableAssets ? 'rotated' : ''}`}
                width={16}
                height={16}
                stroke="currentColor"
                strokeWidth={2}
              />
            )}
          </button>
        </div>
        
        {showAvailableAssets && (
          <div className="asset-list__items">
            {getFilteredAssetTypes().length === 0 ? (
              <div className="asset-list__empty-message">
                No additional assets available
              </div>
            ) : (
              getFilteredAssetTypes().map((asset, index) => (
                <AssetSummary key={`${asset.assetType}-${index}`} asset={asset} />
              ))
            )}
          </div>
        )}
      </div>
    </>
  );
}
