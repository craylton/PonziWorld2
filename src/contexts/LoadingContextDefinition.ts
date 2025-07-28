import { createContext } from 'react';

interface LoadingContextType {
    showLoadingPopup: (status: 'loading' | 'success' | 'error', message?: string, customOnClose?: () => void) => void;
    hideLoadingPopup: () => void;
}

export const LoadingContext = createContext<LoadingContextType | undefined>(undefined);
