export interface User {
  username: string;
}

export interface Bank {
  bankName: string;
  claimedCapital: number;
  actualCapital: number;
  assets: Asset[];
}

export interface Asset {
  amount: number;
  assetType: string;
}
