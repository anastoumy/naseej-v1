import { useState, useCallback } from "react";
import { Search, RotateCcw, Loader2 } from "lucide-react";
import { Modal } from "@/components/ui/custom/modal/Modal";
import { cn, formatSAR } from "@/lib/utils";

interface RefundItem {
  id: string;
  productName: string;
  size: string;
  color: string;
  quantity: number;
  unitPrice: number;
  selected: boolean;
  refundQty: number;
}

interface RefundModalProps {
  open: boolean;
  onClose: () => void;
  onRefund?: (items: RefundItem[], method: string) => void;
}

// Mock order items for refund
const MOCK_ORDER_ITEMS: Omit<RefundItem, "selected" | "refundQty">[] = [
  { id: "1", productName: "Classic T-Shirt", size: "M", color: "Black", quantity: 2, unitPrice: 89.0 },
  { id: "2", productName: "Slim Fit Jeans", size: "32", color: "Dark Blue", quantity: 1, unitPrice: 249.0 },
  { id: "3", productName: "Cotton Hoodie", size: "L", color: "Gray", quantity: 1, unitPrice: 179.0 },
];

export default function RefundModal({
  open,
  onClose,
  onRefund,
}: RefundModalProps) {
  const [receiptNumber, setReceiptNumber] = useState("");
  const [orderLoaded, setOrderLoaded] = useState(false);
  const [items, setItems] = useState<RefundItem[]>([]);
  const [refundMethod, setRefundMethod] = useState<"cash" | "card" | "credit">("cash");
  const [isProcessing, setIsProcessing] = useState(false);

  const handleSearch = useCallback(() => {
    if (!receiptNumber.trim()) return;
    // Simulate loading an order
    setItems(
      MOCK_ORDER_ITEMS.map((item) => ({
        ...item,
        selected: false,
        refundQty: 0,
      })),
    );
    setOrderLoaded(true);
  }, [receiptNumber]);

  const toggleItem = useCallback((itemId: string) => {
    setItems((prev) =>
      prev.map((item) =>
        item.id === itemId
          ? { ...item, selected: !item.selected, refundQty: item.selected ? 0 : item.quantity }
          : item,
      ),
    );
  }, []);

  const selectedItems = items.filter((i) => i.selected);
  const refundTotal = selectedItems.reduce((sum, i) => sum + i.refundQty * i.unitPrice, 0);

  const handleRefund = useCallback(async () => {
    if (selectedItems.length === 0) return;
    setIsProcessing(true);
    await new Promise((resolve) => setTimeout(resolve, 1200));
    setIsProcessing(false);
    onRefund?.(selectedItems, refundMethod);
    onClose();
  }, [selectedItems, refundMethod, onRefund, onClose]);

  const handleClose = useCallback(() => {
    setReceiptNumber("");
    setOrderLoaded(false);
    setItems([]);
    setRefundMethod("cash");
    onClose();
  }, [onClose]);

  const footer = orderLoaded ? (
    <div className="flex items-center justify-between">
      <div>
        <p className="text-sm text-muted-foreground">
          {selectedItems.length} item(s) selected
        </p>
        <p className="text-base font-semibold text-foreground tabular-nums">
          Refund: {formatSAR(refundTotal)}
        </p>
      </div>
      <button
        type="button"
        onClick={handleRefund}
        disabled={isProcessing || selectedItems.length === 0}
        className="btn-primary flex items-center gap-2 px-6 py-2 bg-destructive hover:bg-destructive/90 disabled:opacity-50"
      >
        {isProcessing && <Loader2 className="w-4 h-4 animate-spin" />}
        {isProcessing ? "Processing..." : "Process Refund"}
      </button>
    </div>
  ) : undefined;

  return (
    <Modal
      open={open}
      onClose={handleClose}
      title="Process Refund"
      description="Search receipt to initiate a refund"
      size="lg"
      footer={footer}
    >
      <div className="space-y-4">
        {/* Receipt Search */}
        <div className="flex gap-2">
          <div className="relative flex-1">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
            <input
              type="text"
              value={receiptNumber}
              onChange={(e) => setReceiptNumber(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && handleSearch()}
              placeholder="Enter receipt number or scan barcode..."
              className="pos-input pl-10"
              autoFocus
            />
          </div>
          <button
            onClick={handleSearch}
            className="px-4 py-2 rounded-lg text-sm font-medium bg-accent text-accent-foreground hover:bg-accent/80 transition-colors"
          >
            Search
          </button>
        </div>

        {/* Order Items */}
        {orderLoaded && (
          <>
            <div className="space-y-2">
              {items.map((item) => (
                <label
                  key={item.id}
                  className={cn(
                    "flex items-center gap-3 p-3 rounded-lg border cursor-pointer transition-colors",
                    item.selected
                      ? "border-destructive/50 bg-destructive/5"
                      : "border-border hover:bg-muted/50",
                  )}
                >
                  <input
                    type="checkbox"
                    checked={item.selected}
                    onChange={() => toggleItem(item.id)}
                    className="w-4 h-4 rounded border-border text-destructive focus:ring-destructive"
                  />
                  <div className="flex-1 min-w-0">
                    <p className="text-sm font-medium text-foreground">{item.productName}</p>
                    <p className="text-xs text-muted-foreground">
                      {item.size} · {item.color} · Qty: {item.quantity}
                    </p>
                  </div>
                  <p className="text-sm font-medium text-foreground tabular-nums">
                    {formatSAR(item.unitPrice * item.quantity)}
                  </p>
                </label>
              ))}
            </div>

            {/* Refund Method */}
            <div className="space-y-2">
              <label className="text-sm font-medium text-foreground">Refund Method</label>
              <div className="flex gap-2">
                {(["cash", "card", "credit"] as const).map((method) => (
                  <button
                    key={method}
                    onClick={() => setRefundMethod(method)}
                    className={cn(
                      "px-4 py-2 rounded-lg text-sm font-medium border transition-colors capitalize",
                      refundMethod === method
                        ? "border-accent bg-accent/15 text-accent"
                        : "border-border text-muted-foreground hover:bg-muted",
                    )}
                  >
                    {method === "credit" ? "Store Credit" : method}
                  </button>
                ))}
              </div>
            </div>
          </>
        )}

        {!orderLoaded && (
          <div className="flex flex-col items-center justify-center py-8 text-muted-foreground">
            <RotateCcw className="w-8 h-8 mb-2 opacity-50" />
            <p className="text-sm">Enter a receipt number to load the order</p>
          </div>
        )}
      </div>
    </Modal>
  );
}
