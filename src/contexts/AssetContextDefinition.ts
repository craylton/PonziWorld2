import { createContext } from 'react';

interface AssetContextType {
  refreshAssets: () => void;
  registerRefreshCallback: (callback: () => void) => void;
  unregisterRefreshCallback: (callback: () => void) => void;
  refreshBank?: () => Promise<void>;
  // Cash balance
  cashBalance: number;
  setCashBalance: (balance: number) => void;
}

export const AssetContext = createContext<AssetContextType | undefined>(undefined);
