import { createContext } from 'react';

interface AssetContextType {
  refreshBank?: () => Promise<void>;
  cashBalance: number;
  setCashBalance: (balance: number) => void;
}

export const AssetContext = createContext<AssetContextType | undefined>(undefined);
