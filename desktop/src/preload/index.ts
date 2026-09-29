import { contextBridge, ipcRenderer } from "electron";
import { electronAPI } from "@electron-toolkit/preload";

const api = {
  // Platform detection
  getPlatform: () => process.platform,

  // Window controls
  minimizeWindow: () => ipcRenderer.send("window:minimize"),
  maximizeWindow: () => ipcRenderer.send("window:maximize"),
  closeWindow: () => ipcRenderer.send("window:close"),

  // Hardware
  printReceipt: (data: unknown) =>
    ipcRenderer.invoke("hardware:print-receipt", data),
  openCashDrawer: () => ipcRenderer.invoke("hardware:open-drawer"),

  // Scanner events (main → renderer)
  onBarcodeScanned: (callback: (barcode: string) => void) => {
    const listener = (_event: Electron.IpcRendererEvent, barcode: string) =>
      callback(barcode);
    ipcRenderer.on("scanner:barcode", listener);
    return () => ipcRenderer.removeListener("scanner:barcode", listener);
  },

  // System
  getNetworkStatus: () => ipcRenderer.invoke("system:network-status"),
};

if (process.contextIsolated) {
  try {
    contextBridge.exposeInMainWorld("electron", electronAPI);
    contextBridge.exposeInMainWorld("api", api);
  } catch (error) {
    console.error(error);
  }
} else {
  // @ts-ignore (define in d.ts)
  window.electron = electronAPI;
  // @ts-ignore (define in d.ts)
  window.api = api;
}
