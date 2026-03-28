import { useContext } from 'react';
import { AssetContext } from './AssetContext';

export function useAssetContext() {
  const context = useContext(AssetContext);
  if (context === undefined) {
    throw new Error('useAssetContext must be used within an AssetProvider');
  }
  return context;
}
