import { createContext } from 'react';

interface AssetContextType {
  refreshAssets: () => void;
  registerRefreshCallback: (callback: () => void) => void;
  unregisterRefreshCallback: (callback: () => void) => void;
  refreshBank?: () => Promise<void>;
  // Cash balance - now string for arbitrary precision
  cashBalance: string;
  setCashBalance: (balance: string) => void;
}

export const AssetContext = createContext<AssetContextType | undefined>(undefined);
