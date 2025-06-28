export interface PerformanceHistoryEntry {
    day: number;
    value: number;
}

export interface PerformanceHistory {
    claimedHistory: PerformanceHistoryEntry[];
    actualHistory: PerformanceHistoryEntry[];
}
