import { useState, useCallback } from "react";
import { Loader2 } from "lucide-react";
import { Modal } from "@/components/ui/custom/modal/Modal";
import { mockCategories, mockSizes } from "../__mocks__/categoryMocks";
import { CategoryInfoFields } from "../components/CategoryInfoFields";
import { StandardSizesToggle } from "../components/StandardSizesToggle";

interface EditCategoryModalProps {
  open: boolean;
  onClose: () => void;
  onSuccess?: () => void;
  category?: {
    id: string;
    name: string;
    slug: string;
    parentId: string | null;
    hasStandardSizes: boolean;
    selectedSizes: string[];
  } | null;
}

export default function EditCategoryModal({
  open,
  onClose,
  onSuccess,
  category,
}: EditCategoryModalProps) {
  const [name, setName] = useState(category?.name ?? "");
  const [slug, setSlug] = useState(category?.slug ?? "");
  const [parentId, setParentId] = useState(category?.parentId ?? null);
  const [hasStandardSizes, setHasStandardSizes] = useState(category?.hasStandardSizes ?? false);
  const [selectedSizes, setSelectedSizes] = useState<string[]>(category?.selectedSizes ?? []);
  const [isSaving, setIsSaving] = useState(false);
  const [errors, setErrors] = useState<Record<string, string>>({});

  const handleClose = useCallback(() => {
    setErrors({});
    onClose();
  }, [onClose]);

  const handleSave = useCallback(async () => {
    const newErrors: Record<string, string> = {};
    if (!name.trim()) newErrors.name = "Category name is required";
    if (!slug.trim()) newErrors.slug = "Slug is required";

    if (Object.keys(newErrors).length > 0) {
      setErrors(newErrors);
      return;
    }

    setIsSaving(true);
    await new Promise((resolve) => setTimeout(resolve, 800));
    setIsSaving(false);
    onSuccess?.();
    handleClose();
  }, [name, slug, onSuccess, handleClose]);

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
        {isSaving ? "Saving..." : "Update Category"}
      </button>
    </div>
  );

  return (
    <Modal
      open={open}
      onClose={handleClose}
      title="Edit Category"
      description="Modify category details and size configuration"
      size="lg"
      footer={footer}
    >
      <form onSubmit={(e) => e.preventDefault()} className="space-y-6">
        <section className="space-y-4">
          <CategoryInfoFields
            name={name}
            onNameChange={(val) => setName(val)}
            slug={slug}
            onSlugChange={(val) => setSlug(val)}
            parentId={parentId}
            onParentIdChange={(val) => setParentId(val)}
            categories={mockCategories.filter((c) => c.id !== category?.id)}
            isParentLocked={false}
            nameError={errors.name}
            slugError={errors.slug}
          />
        </section>

        <section className="space-y-4">
          <StandardSizesToggle
            hasStandardSizes={hasStandardSizes}
            onToggle={() => setHasStandardSizes(!hasStandardSizes)}
            sizes={mockSizes}
            selectedSizes={selectedSizes}
            onToggleSize={(sizeId) =>
              setSelectedSizes((prev) =>
                prev.includes(sizeId) ? prev.filter((s) => s !== sizeId) : [...prev, sizeId]
              )
            }
            onSelectAll={() => setSelectedSizes(mockSizes.map((s) => s.id))}
            onClearAll={() => setSelectedSizes([])}
            error={errors.selectedSizes}
          />
        </section>
      </form>
    </Modal>
  );
}
