import type { Asset } from "./Asset";


export interface Bank {
    id: string;
    bankName: string;
    claimedCapital: number;
    actualCapital: number;
    assets: Asset[];
}
