import type { AvailableAsset } from "./AvailableAsset";
import type { Investor } from "../Dashboard/SidePanel/Investors/Investor";

export interface Bank {
    id: string;
    bankName: string;
    claimedCapital: number;
    actualCapital: number;
    availableAssets: AvailableAsset[];
    investors: Investor[];
}
