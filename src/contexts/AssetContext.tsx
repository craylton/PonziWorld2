import { useState } from 'react';
import { AssetContext } from './AssetContextDefinition';

interface AssetProviderProps {
  children: React.ReactNode;
  refreshBank?: () => Promise<void>;
}

export default function AssetProvider({ children, refreshBank }: AssetProviderProps) {
  const [cashBalance, setCashBalance] = useState<number>(0);

  const value = {
    refreshBank,
    cashBalance,
    setCashBalance,
  };

  return (
    <AssetContext.Provider value={value}>
      {children}
    </AssetContext.Provider>
  );
}
