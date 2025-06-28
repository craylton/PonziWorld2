import type { Asset } from './Dashboard/AssetList/Asset';

export interface User {
  username: string;
}

export interface Bank {
  id: string;
  bankName: string;
  claimedCapital: number;
  actualCapital: number;
  assets: Asset[];
}

export interface PerformanceHistoryEntry {
  day: number;
  value: number;
}

export interface PerformanceHistory {
  claimedHistory: PerformanceHistoryEntry[];
  actualHistory?: PerformanceHistoryEntry[];
}
