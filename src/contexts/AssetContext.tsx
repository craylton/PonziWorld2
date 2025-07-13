import React, { useCallback } from 'react';
import { AssetContext } from './AssetContextDefinition';

interface AssetProviderProps {
  children: React.ReactNode;
}

export default function AssetProvider({ children }: AssetProviderProps) {
  const refreshCallbacks = React.useRef<Set<() => void>>(new Set());

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
  };

  return (
    <AssetContext.Provider value={value}>
      {children}
    </AssetContext.Provider>
  );
}
