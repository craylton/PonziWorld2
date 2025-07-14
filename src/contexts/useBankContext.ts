import { useContext } from 'react';
import { BankContext } from './BankContext';
import type { BankContextType } from './BankContext';

export function useBankContext(): BankContextType {
  const context = useContext(BankContext);
  if (context === undefined) {
    throw new Error('useBankContext must be used within a BankProvider');
  }
  return context;
}
