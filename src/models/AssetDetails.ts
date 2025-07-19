import type { HistoricalPerformanceEntry } from './HistoricalPerformance';

export interface AssetDetailsResponse {
  assetId: string;
  name: string;
  investedAmount: number;
  pendingAmount: number;
  historicalData: HistoricalPerformanceEntry[];
}
