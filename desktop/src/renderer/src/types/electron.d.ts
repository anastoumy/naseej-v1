// Type augmentation for Electron preload API
interface ElectronAPI {
  getPlatform: () => string;
  minimizeWindow: () => void;
  maximizeWindow: () => void;
  closeWindow: () => void;
  printReceipt: (data: unknown) => Promise<{ success: boolean }>;
  openCashDrawer: () => Promise<{ success: boolean }>;
  onBarcodeScanned: (callback: (barcode: string) => void) => () => void;
  getNetworkStatus: () => Promise<{ online: boolean }>;
}

declare global {
  interface Window {
    api?: ElectronAPI;
    electron?: unknown;
  }
}

export {};
