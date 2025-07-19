import type { HistoricalPerformanceEntry } from './HistoricalPerformance';

export interface InvestmentDetailsResponse {
  targetAssetId: string;
  name: string;
  investedAmount: number;
  pendingAmount: number;
  historicalData: HistoricalPerformanceEntry[];
}
