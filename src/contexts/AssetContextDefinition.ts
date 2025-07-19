import { createContext } from 'react';

interface AssetContextType {
  refreshAssets: () => void;
  registerRefreshCallback: (callback: () => void) => void;
  unregisterRefreshCallback: (callback: () => void) => void;
  refreshBank?: () => Promise<void>;
  // Global loading popup state
  showLoadingPopup: (status: 'loading' | 'success' | 'error', message?: string) => void;
  hideLoadingPopup: () => void;
}

export const AssetContext = createContext<AssetContextType | undefined>(undefined);
