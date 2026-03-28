import { useCallback, useState } from 'react';
import type { ReactNode } from 'react';

import LoadingPopup from '../Dashboard/Assets/LoadingPopup';
import { LoadingContext } from './LoadingContext';

interface LoadingProviderProps {
  children: ReactNode;
}

export function LoadingProvider({ children }: LoadingProviderProps) {
  const [loadingPopupOpen, setLoadingPopupOpen] = useState(false);
  const [loadingStatus, setLoadingStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [loadingMessage, setLoadingMessage] = useState<string>('');
  const [customOnClose, setCustomOnClose] = useState<(() => void) | undefined>();

  const showLoadingPopup = useCallback((status: 'loading' | 'success' | 'error', message?: string, customCloseHandler?: () => void) => {
    setLoadingStatus(status);
    setLoadingMessage(message || '');
    setCustomOnClose(() => customCloseHandler);
    setLoadingPopupOpen(true);
  }, []);

  const hideLoadingPopup = useCallback(() => {
    setLoadingPopupOpen(false);
    if (customOnClose) {
      customOnClose();
      setCustomOnClose(undefined);
    }
  }, [customOnClose]);

  const value = {
    showLoadingPopup,
    hideLoadingPopup,
  };

  return (
    <LoadingContext.Provider value={value}>
      {children}
      <LoadingPopup
        isOpen={loadingPopupOpen}
        onClose={hideLoadingPopup}
        status={loadingStatus}
        message={loadingMessage}
      />
    </LoadingContext.Provider>
  );
}