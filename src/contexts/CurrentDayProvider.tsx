import { useEffect, useState } from 'react';
import type { ReactNode } from 'react';

import { CurrentDayContext } from './currentDayContext';

export function CurrentDayProvider({ children }: { children: ReactNode }) {
  const [currentDay, setCurrentDay] = useState<number | null>(null);

  const fetchCurrentDay = async () => {
    const response = await fetch('/api/currentDay');
    if (response.ok) {
      const data: { currentDay: number } = await response.json();
      setCurrentDay(data.currentDay);
    }
  };

  useEffect(() => {
    fetchCurrentDay();
  }, []);

  return (
    <CurrentDayContext.Provider value={{ currentDay, refreshCurrentDay: fetchCurrentDay }}>
      {children}
    </CurrentDayContext.Provider>
  );
}
