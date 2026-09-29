import { useState } from "react";
import {
  Search,
  Package,
  MapPin,
  TrendingDown,
  AlertTriangle,
  ChevronDown,
  ChevronUp,
  Plus,
  Pencil,
} from "lucide-react";
import { useTranslation } from "react-i18next";
import type { InventoryItem } from "@/types/pos";
import { DUMMY_INVENTORY, DUMMY_PRODUCTS } from "@/data/dummy-data";
import { cn } from "@/lib/utils";
import CreateProductModal from "../modals/CreateProductModal";
import EditProductModal from "../modals/EditProductModal";

export default function InventoryPage() {
  const { t } = useTranslation();
  const [searchQuery, setSearchQuery] = useState("");
  const [expandedItem, setExpandedItem] = useState<string | null>(null);
  const [createOpen, setCreateOpen] = useState(false);
  const [editProduct, setEditProduct] = useState<{
    id: string;
    name: string;
    modelNumber: string;
    basePrice: number;
    brandId: string;
    categoryId: string;
    description: string;
  } | null>(null);

  const handleEditClick = (item: InventoryItem, e: React.MouseEvent) => {
    e.stopPropagation();
    const product = DUMMY_PRODUCTS.find((p) => p.id === item.productId);
    setEditProduct({
      id: item.productId,
      name: item.productName,
      modelNumber: product?.sku ?? "",
      basePrice: item.price,
      brandId: "",
      categoryId: "",
      description: "",
    });
  };

  const filtered = DUMMY_INVENTORY.filter(
    (item) =>
      item.productName.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.sku.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.barcode.includes(searchQuery) ||
      item.size.toLowerCase().includes(searchQuery.toLowerCase()) ||
      item.color.toLowerCase().includes(searchQuery.toLowerCase()),
  );

  const getTotalQty = (item: InventoryItem) =>
    item.locations.reduce((sum, loc) => sum + loc.qty, 0);

  const isLowStock = (item: InventoryItem) =>
    item.locations.some(
      (loc) => loc.qty > 0 && loc.qty <= loc.lowStockThreshold,
    );

  const isOutOfStock = (item: InventoryItem) =>
    item.locations.every((loc) => loc.qty === 0);

  return (
    <div className="flex-1 flex flex-col overflow-hidden">
      {/* Header */}
      <div className="h-12 bg-card border-b border-border flex items-center gap-3 px-4 shrink-0">
        <Package className="w-4 h-4 text-muted-foreground" />
        <h1 className="text-sm font-semibold text-foreground">
          {t('features.catalog.screens.inventory.title')}
        </h1>
        <div className="ml-auto flex items-center gap-2 flex-1 max-w-sm">
          <div className="relative flex-1">
            <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-muted-foreground" />
            <input
              type="text"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              placeholder={t('features.catalog.screens.inventory.searchPlaceholder')}
              className="pos-input pl-8 text-sm h-8"
              aria-label="Search inventory"
            />
          </div>
        </div>
        <button
          type="button"
          onClick={() => setCreateOpen(true)}
          className="btn-primary flex items-center gap-1.5 px-3 py-1.5 text-sm shrink-0"
        >
          <Plus className="w-3.5 h-3.5" />
          Add Product
        </button>
      </div>

      {/* Stats bar */}
      <div className="flex gap-px bg-border shrink-0">
        {[
          {
            label: t('features.catalog.screens.inventory.totalSKUs'),
            value: DUMMY_INVENTORY.length,
            color: "text-foreground",
          },
          {
            label: t('features.catalog.screens.inventory.lowStock'),
            value: DUMMY_INVENTORY.filter(isLowStock).length,
            color: "text-warning",
          },
          {
            label: t('features.catalog.screens.inventory.outOfStock'),
            value: DUMMY_INVENTORY.filter(isOutOfStock).length,
            color: "text-destructive",
          },
        ].map(({ label, value, color }) => (
          <div key={label} className="flex-1 bg-card px-4 py-2">
            <p className={cn("text-lg font-bold tabular-nums", color)}>
              {value}
            </p>
            <p className="text-xs text-muted-foreground">{label}</p>
          </div>
        ))}
      </div>

      {/* Table */}
      <div className="flex-1 overflow-y-auto">
        {/* Table header */}
        <div className="grid grid-cols-[2fr_1fr_1fr_1fr_auto] gap-0 bg-muted border-b border-border px-4 py-2 sticky top-0 z-10">
          {[t('features.catalog.screens.inventory.productSKU'), t('features.catalog.screens.inventory.size'), t('features.catalog.screens.inventory.color'), t('features.catalog.screens.inventory.totalStock'), t('features.catalog.screens.inventory.details')].map(
            (col) => (
              <div
                key={col}
                className="text-xs font-semibold text-muted-foreground uppercase tracking-wide"
              >
                {col}
              </div>
            ),
          )}
        </div>

        {filtered.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 text-center">
            <Package className="w-10 h-10 text-muted-foreground/30 mb-3" />
            <p className="text-sm font-medium text-foreground">
              {t('features.catalog.screens.inventory.noVariantsFound')}
            </p>
            <p className="text-xs text-muted-foreground mt-1">
              {t('features.catalog.screens.inventory.tryDifferentSearch')}
            </p>
          </div>
        ) : (
          filtered.map((item) => {
            const totalQty = getTotalQty(item);
            const lowStock = isLowStock(item);
            const outOfStock = isOutOfStock(item);
            const isExpanded = expandedItem === item.variantId;

            return (
              <div key={item.variantId} className="border-b border-border/50">
                <div
                  className={cn(
                    "grid grid-cols-[2fr_1fr_1fr_1fr_auto] gap-0 px-4 py-3 items-center",
                    "hover:bg-muted/50 transition-colors cursor-pointer",
                    outOfStock && "opacity-60",
                  )}
                  onClick={() =>
                    setExpandedItem(isExpanded ? null : item.variantId)
                  }
                  role="button"
                  tabIndex={0}
                  onKeyDown={(e) =>
                    e.key === "Enter" &&
                    setExpandedItem(isExpanded ? null : item.variantId)
                  }
                  aria-expanded={isExpanded}
                >
                  {/* Product name + SKU */}
                  <div>
                    <p className="text-sm font-medium text-foreground">
                      {item.productName}
                    </p>
                    <p className="text-xs text-muted-foreground font-mono">
                      {item.sku}
                    </p>
                  </div>

                  {/* Size */}
                  <span className="badge-muted w-fit">{item.size}</span>

                  {/* Color */}
                  <div className="flex items-center gap-2">
                    <div
                      className="w-4 h-4 rounded border border-border/50"
                      style={{ backgroundColor: item.colorHex }}
                    />
                    <span className="text-sm text-foreground">
                      {item.color}
                    </span>
                  </div>

                  {/* Total stock + status */}
                  <div className="flex items-center gap-2">
                    <span
                      className={cn(
                        "text-sm font-semibold tabular-nums",
                        outOfStock
                          ? "text-destructive"
                          : lowStock
                            ? "text-warning"
                            : "text-success",
                      )}
                    >
                      {totalQty}
                    </span>
                    {outOfStock && (
                      <AlertTriangle className="w-3.5 h-3.5 text-destructive" />
                    )}
                    {lowStock && !outOfStock && (
                      <TrendingDown className="w-3.5 h-3.5 text-warning" />
                    )}
                    {outOfStock ? (
                      <span className="badge-destructive">Out</span>
                    ) : lowStock ? (
                      <span className="badge-warning">Low</span>
                    ) : null}
                  </div>

                  {/* Actions */}
                  <div className="flex items-center gap-1">
                    <button
                      className="w-6 h-6 flex items-center justify-center text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
                      aria-label="Edit product"
                      onClick={(e) => handleEditClick(item, e)}
                    >
                      <Pencil className="w-3.5 h-3.5" />
                    </button>
                    <button
                      className="w-6 h-6 flex items-center justify-center text-muted-foreground cursor-pointer"
                      aria-label={isExpanded ? "Collapse" : "Expand locations"}
                    >
                      {isExpanded ? (
                        <ChevronUp className="w-4 h-4" />
                      ) : (
                        <ChevronDown className="w-4 h-4" />
                      )}
                    </button>
                  </div>
                </div>

                {/* Expanded: location breakdown */}
                {isExpanded && (
                  <div className="bg-muted/30 border-t border-border/30 px-8 py-3 space-y-2 animate-fade-in">
                    {item.locations.map((loc) => (
                      <div
                        key={loc.locationId}
                        className="flex items-center justify-between"
                      >
                        <div className="flex items-center gap-2">
                          <MapPin className="w-3.5 h-3.5 text-muted-foreground" />
                          <span className="text-sm text-foreground">
                            {loc.locationName}
                          </span>
                        </div>
                        <div className="flex items-center gap-3">
                          <span
                            className={cn(
                              "text-sm font-semibold tabular-nums",
                              loc.qty === 0
                                ? "text-destructive"
                                : loc.qty <= loc.lowStockThreshold
                                  ? "text-warning"
                                  : "text-success",
                            )}
                          >
                            {loc.qty} units
                          </span>
                          {loc.qty === 0 && (
                            <span className="badge-destructive">Out</span>
                          )}
                          {loc.qty > 0 && loc.qty <= loc.lowStockThreshold && (
                            <span className="badge-warning">
                              Low ≤{loc.lowStockThreshold}
                            </span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                )}
              </div>
            );
          })
        )}
      </div>
      <CreateProductModal
        open={createOpen}
        onClose={() => setCreateOpen(false)}
        onSuccess={() => setCreateOpen(false)}
      />
      <EditProductModal
        open={editProduct !== null}
        onClose={() => setEditProduct(null)}
        onSuccess={() => setEditProduct(null)}
        product={editProduct}
      />
    </div>
  );
}
