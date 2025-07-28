import { useContext } from 'react';
import { LoadingContext } from './LoadingContextDefinition';

export function useLoadingContext() {
    const context = useContext(LoadingContext);
    if (context === undefined) {
        throw new Error('useLoadingContext must be used within a LoadingProvider');
    }
    return context;
}
