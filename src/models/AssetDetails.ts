import type { HistoricalPerformanceEntry } from './HistoricalPerformance';

export interface InvestmentDetailsResponse {
  targetAssetId: string;
  targetAssetName: string;
  investedAmount: number;
  pendingAmount: number;
  historicalData: HistoricalPerformanceEntry[];
}
