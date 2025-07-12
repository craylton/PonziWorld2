export interface PendingTransaction {
  id: string;
  buyerBankId: string;
  assetId: string;
  amount: number; // Positive = buy, negative = sell
  createdAt: string;
}
