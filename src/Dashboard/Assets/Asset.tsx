export interface Asset {
  amount: number;
  assetType: string;
  assetTypeId: string;
  pendingAmount: number; // Pending transaction amount (positive for buy, negative for sell)
}
