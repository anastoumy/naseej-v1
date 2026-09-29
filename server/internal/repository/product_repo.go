package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/pkg/utils"
)

// ProductRepository handles product data access.
type ProductRepository struct {
	*BaseRepository
}

// NewProductRepository creates a new ProductRepository.
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{BaseRepository: NewBaseRepository(db)}
}

// FindByID retrieves a product by ID with all associations.
func (r *ProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.Scoped(ctx).
		Preload("Brand").
		Preload("Category").
		Preload("Images").
		Preload("Variants").
		Preload("Variants.Size").
		Preload("Variants.Color").
		Where("id = ?", id).
		First(&product).Error; err != nil {
		return nil, fmt.Errorf("find product by id: %w", err)
	}
	return &product, nil
}

// FindBySlug retrieves a product by slug.
func (r *ProductRepository) FindBySlug(ctx context.Context, slug string) (*models.Product, error) {
	var product models.Product
	if err := r.Scoped(ctx).
		Preload("Brand").
		Preload("Category").
		Preload("Images").
		Preload("Variants").
		Preload("Variants.Size").
		Preload("Variants.Color").
		Where("slug = ?", slug).
		First(&product).Error; err != nil {
		return nil, fmt.Errorf("find product by slug: %w", err)
	}
	return &product, nil
}

// List retrieves products with pagination.
func (r *ProductRepository) List(ctx context.Context, pagination utils.PaginationInput) ([]models.Product, int64, error) {
	var products []models.Product
	var total int64

	query := r.Scoped(ctx).Model(&models.Product{})
	query.Count(&total)

	if err := query.
		Preload("Brand").
		Preload("Category").
		Preload("Images").
		Preload("Variants").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order("created_at DESC").
		Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list products: %w", err)
	}
	return products, total, nil
}

// Create creates a new product.
func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	if err := r.DB.Create(product).Error; err != nil {
		return fmt.Errorf("create product: %w", err)
	}
	return nil
}

// Update updates an existing product.
func (r *ProductRepository) Update(ctx context.Context, product *models.Product) error {
	if err := r.Scoped(ctx).Save(product).Error; err != nil {
		return fmt.Errorf("update product: %w", err)
	}
	return nil
}

// FindVariantByBarcode retrieves a product variant by barcode.
func (r *ProductRepository) FindVariantByBarcode(ctx context.Context, barcode string) (*models.ProductVariant, error) {
	var variant models.ProductVariant
	if err := r.DB.
		Preload("Product").
		Preload("Size").
		Preload("Color").
		Where("barcode = ?", barcode).
		First(&variant).Error; err != nil {
		return nil, fmt.Errorf("find variant by barcode: %w", err)
	}
	return &variant, nil
}

// CreateVariant creates a new product variant.
func (r *ProductRepository) CreateVariant(ctx context.Context, variant *models.ProductVariant) error {
	if err := r.DB.Create(variant).Error; err != nil {
		return fmt.Errorf("create variant: %w", err)
	}
	return nil
}

// ListCategories retrieves all categories for the merchant.
func (r *ProductRepository) ListCategories(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	if err := r.Scoped(ctx).
		Preload("SizeGroup").
		Preload("SizeGroup.Sizes").
		Preload("Children").
		Where("parent_id IS NULL").
		Order("sort_order").
		Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}

// ListSizeGroups retrieves all size groups for the merchant.
func (r *ProductRepository) ListSizeGroups(ctx context.Context) ([]models.SizeGroup, error) {
	var groups []models.SizeGroup
	if err := r.Scoped(ctx).
		Preload("Sizes", func(db *gorm.DB) *gorm.DB {
			return db.Order("sort_order")
		}).
		Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("list size groups: %w", err)
	}
	return groups, nil
}

// CreateSizeGroup creates a new size group with sizes.
func (r *ProductRepository) CreateSizeGroup(ctx context.Context, group *models.SizeGroup) error {
	if err := r.DB.Create(group).Error; err != nil {
		return fmt.Errorf("create size group: %w", err)
	}
	return nil
}
