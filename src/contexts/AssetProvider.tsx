import { useState } from 'react';
import type { ReactNode } from 'react';

import { AssetContext } from './AssetContext';

interface AssetProviderProps {
  children: ReactNode;
  refreshBank?: () => Promise<void>;
}

export function AssetProvider({ children, refreshBank }: AssetProviderProps) {
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