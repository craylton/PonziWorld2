import React, { useCallback, useState } from 'react';
import { AssetContext } from './AssetContextDefinition';
import LoadingPopup from '../Dashboard/Assets/LoadingPopup';

interface AssetProviderProps {
  children: React.ReactNode;
  refreshBank?: () => Promise<void>;
}

export default function AssetProvider({ children, refreshBank }: AssetProviderProps) {
  const refreshCallbacks = React.useRef<Set<() => void>>(new Set());
  
  // Global loading popup state
  const [loadingPopupOpen, setLoadingPopupOpen] = useState(false);
  const [loadingStatus, setLoadingStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [loadingMessage, setLoadingMessage] = useState<string>('');

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

  const showLoadingPopup = useCallback((status: 'loading' | 'success' | 'error', message?: string) => {
    setLoadingStatus(status);
    setLoadingMessage(message || '');
    setLoadingPopupOpen(true);
  }, []);

  const hideLoadingPopup = useCallback(() => {
    setLoadingPopupOpen(false);
  }, []);

  const value = {
    refreshAssets,
    registerRefreshCallback,
    unregisterRefreshCallback,
    refreshBank,
    showLoadingPopup,
    hideLoadingPopup,
  };

  return (
    <AssetContext.Provider value={value}>
      {children}
      <LoadingPopup
        isOpen={loadingPopupOpen}
        onClose={hideLoadingPopup}
        status={loadingStatus}
        message={loadingMessage}
      />
    </AssetContext.Provider>
  );
}
