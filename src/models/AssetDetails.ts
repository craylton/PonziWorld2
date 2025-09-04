import type { HistoricalPerformanceEntry } from './HistoricalPerformance';

export interface InvestmentDetailsResponse {
  targetAssetId: string;
  targetAssetName: string;
  investedAmount: string; // Now string for arbitrary precision
  pendingAmount: string;  // Now string for arbitrary precision
  historicalData: HistoricalPerformanceEntry[];
}
