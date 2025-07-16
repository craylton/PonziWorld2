export interface HistoricalPerformanceEntry {
    day: number;
    value: number;
}

export interface HistoricalPerformance {
    claimedHistory: HistoricalPerformanceEntry[];
    actualHistory: HistoricalPerformanceEntry[];
}
