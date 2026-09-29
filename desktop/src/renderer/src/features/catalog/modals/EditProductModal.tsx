import { useState, useCallback } from "react";
import { Loader2 } from "lucide-react";
import { Modal } from "@/components/ui/custom/modal/Modal";
import SearchableCombobox from "../components/SearchableCombobox";
import { mockBrands, mockCategories } from "../__mocks__/productMocks";

interface EditProductModalProps {
  open: boolean;
  onClose: () => void;
  onSuccess?: () => void;
  product?: {
    id: string;
    name: string;
    modelNumber: string;
    basePrice: number;
    brandId: string;
    categoryId: string;
    description: string;
  } | null;
}

export default function EditProductModal({
  open,
  onClose,
  onSuccess,
  product,
}: EditProductModalProps) {
  const [name, setName] = useState(product?.name ?? "");
  const [modelNumber, setModelNumber] = useState(product?.modelNumber ?? "");
  const [basePrice, setBasePrice] = useState(String(product?.basePrice ?? ""));
  const [brandId, setBrandId] = useState(product?.brandId ?? "");
  const [categoryId, setCategoryId] = useState(product?.categoryId ?? "");
  const [description, setDescription] = useState(product?.description ?? "");
  const [isSaving, setIsSaving] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleClose = useCallback(() => {
    setErrors({});
    onClose();
  }, [onClose]);

  const handleSave = useCallback(async () => {
    const newErrors: Record<string, string> = {};
    if (!name.trim()) newErrors.name = "Product name is required";
    if (!modelNumber.trim()) newErrors.modelNumber = "Model number is required";
    if (!basePrice || parseFloat(basePrice) <= 0) newErrors.basePrice = "Valid price is required";

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setIsSaving(true);
    await new Promise((resolve) => setTimeout(resolve, 800));
    setIsSaving(false);
    onSuccess?.();
    handleClose();
  }, [name, modelNumber, basePrice, onSuccess, handleClose]);

  const footer = (
    <div className="flex justify-end gap-3">
      <button
        type="button"
        onClick={handleClose}
        className="px-4 py-2 rounded-lg text-sm font-medium bg-muted text-muted-foreground border border-border hover:bg-muted/80 transition-colors"
      >
        Cancel
      </button>
      <button
        type="button"
        onClick={handleSave}
        disabled={isSaving}
        className="btn-primary flex items-center gap-2 px-6 py-2"
      >
        {isSaving && <Loader2 className="w-4 h-4 animate-spin" />}
        {isSaving ? "Saving..." : "Update Product"}
      </button>
    </div>
  );

  return (
    <Modal
      open={open}
      onClose={handleClose}
      title="Edit Product"
      description="Modify product details"
      size="lg"
      footer={footer}
    >
      <form onSubmit={(e) => e.preventDefault()} className="space-y-4">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">
              Product Name <span className="text-destructive">*</span>
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className={`pos-input ${errors.name ? "border-destructive ring-1 ring-destructive/30" : ""}`}
            />
            {errors.name && <p className="text-xs text-destructive">{errors.name}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">
              Model Number <span className="text-destructive">*</span>
            </label>
            <input
              type="text"
              value={modelNumber}
              onChange={(e) => setModelNumber(e.target.value)}
              className={`pos-input ${errors.modelNumber ? "border-destructive ring-1 ring-destructive/30" : ""}`}
            />
            {errors.modelNumber && <p className="text-xs text-destructive">{errors.modelNumber}</p>}
          </div>

          <div className="space-y-1.5">
            <label className="text-sm font-medium text-foreground">
              Price (SAR) <span className="text-destructive">*</span>
            </label>
            <input
              type="number"
              inputMode="decimal"
              min="0"
              step="0.01"
              value={basePrice}
              onChange={(e) => setBasePrice(e.target.value)}
              className={`pos-input tabular-nums ${errors.basePrice ? "border-destructive ring-1 ring-destructive/30" : ""}`}
            />
            {errors.basePrice && <p className="text-xs text-destructive">{errors.basePrice}</p>}
          </div>

          <SearchableCombobox
            label="Brand"
            placeholder="Select brand"
            options={mockBrands}
            value={brandId}
            onChange={setBrandId}
          />

          <SearchableCombobox
            label="Category"
            placeholder="Select category"
            options={mockCategories}
            value={categoryId}
            onChange={setCategoryId}
            showHierarchy
          />
        </div>

        <div className="space-y-1.5">
          <label className="text-sm font-medium text-foreground">Description</label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={2}
            className="pos-input resize-none min-h-[60px]"
          />
        </div>
      </form>
    </Modal>
  );
}
