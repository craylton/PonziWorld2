import { createContext } from 'react';

export interface AssetContextType {
  refreshBank?: () => Promise<void>;
  cashBalance: number;
  setCashBalance: (balance: number) => void;
}

export const AssetContext = createContext<AssetContextType | undefined>(undefined);