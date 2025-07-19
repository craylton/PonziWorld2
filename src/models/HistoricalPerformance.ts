export interface HistoricalPerformanceEntry {
    day: number;
    value: number;
}

export interface OwnBankHistoricalPerformance {
    claimedHistory: HistoricalPerformanceEntry[];
    actualHistory: HistoricalPerformanceEntry[];
}
