import { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';

interface CurrentDayContextType {
    currentDay: number | null;
    refreshCurrentDay: () => Promise<void>;
}

const CurrentDayContext = createContext<CurrentDayContextType | undefined>(undefined);

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

export function useCurrentDay() {
    const context = useContext(CurrentDayContext);
    if (context === undefined) {
        throw new Error('useCurrentDay must be used within a CurrentDayProvider');
    }
    return context;
}

// Export the context for potential direct usage
export { CurrentDayContext };
