import { createContext } from 'react';

export interface CurrentDayContextType {
  currentDay: number | null;
  refreshCurrentDay: () => Promise<void>;
}

export const CurrentDayContext = createContext<CurrentDayContextType | undefined>(undefined);
