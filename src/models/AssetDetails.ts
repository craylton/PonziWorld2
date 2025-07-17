import type { HistoricalPerformanceEntry } from './HistoricalPerformance';

export interface AssetDetailsResponse {
  investedAmount: number;
  pendingAmount: number;
  historicalData: HistoricalPerformanceEntry[];
}
