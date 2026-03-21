import { useContext } from 'react';

import { CurrentDayContext } from './currentDayContext';

export function useCurrentDay() {
  const context = useContext(CurrentDayContext);
  if (context === undefined) {
    throw new Error('useCurrentDay must be used within a CurrentDayProvider');
  }
  return context;
}
