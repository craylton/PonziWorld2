import type { ReactNode } from 'react';

import { BankContext } from './BankContext';

interface BankProviderProps {
  children: ReactNode;
  bankId: string;
}

export function BankProvider({ children, bankId }: BankProviderProps) {
  return (
    <BankContext.Provider value={{ bankId }}>
      {children}
    </BankContext.Provider>
  );
}