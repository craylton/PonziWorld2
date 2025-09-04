import React, { useCallback, useState } from 'react';
import { AssetContext } from './AssetContextDefinition';

interface AssetProviderProps {
  children: React.ReactNode;
  refreshBank?: () => Promise<void>;
}

export default function AssetProvider({ children, refreshBank }: AssetProviderProps) {
  const refreshCallbacks = React.useRef<Set<() => void>>(new Set());

  // Cash balance state - now string for arbitrary precision
  const [cashBalance, setCashBalance] = useState<string>('0');

  const registerRefreshCallback = useCallback((callback: () => void) => {
    refreshCallbacks.current.add(callback);
  }, []);

  const unregisterRefreshCallback = useCallback((callback: () => void) => {
    refreshCallbacks.current.delete(callback);
  }, []);

  const refreshAssets = useCallback(() => {
    // Trigger all registered refresh callbacks
    refreshCallbacks.current.forEach(callback => {
      try {
        callback();
      } catch (error) {
        console.error('Error in asset refresh callback:', error);
      }
    });
  }, []);

  const value = {
    refreshAssets,
    registerRefreshCallback,
    unregisterRefreshCallback,
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
