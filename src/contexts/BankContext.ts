import { createContext } from 'react';

export interface BankContextType {
  bankId: string;
}

export const BankContext = createContext<BankContextType | undefined>(undefined);