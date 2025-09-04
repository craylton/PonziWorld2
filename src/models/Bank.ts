import type { AvailableAsset } from "./AvailableAsset";
import type { Investor } from "../Dashboard/SidePanel/Investors/Investor";

export interface Bank {
    id: string;
    bankName: string;
    claimedCapital: string; // Now string for arbitrary precision
    actualCapital: string;  // Now string for arbitrary precision
    availableAssets: AvailableAsset[];
    investors: Investor[];
}
