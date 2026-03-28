import { useContext } from 'react';

import { CurrentDayContext } from './currentDayContext';

export function useCurrentDayContext() {
  const context = useContext(CurrentDayContext);
  if (context === undefined) {
    throw new Error('useCurrentDayContext must be used within a CurrentDayProvider');
  }
  return context;
}