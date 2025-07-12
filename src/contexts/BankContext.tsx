import { createContext } from 'react';
import type { ReactNode } from 'react';

interface BankContextType {
  bankId: string;
}

const BankContext = createContext<BankContextType | undefined>(undefined);

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

export { BankContext };
export type { BankContextType };
