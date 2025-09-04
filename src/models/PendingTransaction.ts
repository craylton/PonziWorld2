export interface PendingTransaction {
  id: string;
  buyerBankId: string;
  assetId: string;
  amount: string; // Now string for arbitrary precision
}
