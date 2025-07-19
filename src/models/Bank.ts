import type { AvailableAsset } from "./AvailableAsset";

export interface Bank {
    id: string;
    bankName: string;
    claimedCapital: number;
    actualCapital: number;
    availableAssets: AvailableAsset[];
}
