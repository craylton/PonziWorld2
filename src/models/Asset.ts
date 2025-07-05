export interface Asset {
    amount: number;
    assetTypeId: string; // Using string ID reference to AssetType
    assetType: string; // Asset type name for convenience
}
