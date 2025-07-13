import { createContext } from 'react';

interface AssetContextType {
  refreshAssets: () => void;
  registerRefreshCallback: (callback: () => void) => void;
  unregisterRefreshCallback: (callback: () => void) => void;
}

export const AssetContext = createContext<AssetContextType | undefined>(undefined);
