package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/naseej/server/internal/middleware"
	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/repository"
	"github.com/naseej/server/pkg/utils"
)

// ProductService handles product business logic.
type ProductService struct {
	productRepo *repository.ProductRepository
}

// NewProductService creates a new ProductService.
func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

// GetProduct retrieves a product by slug.
func (s *ProductService) GetProduct(ctx context.Context, slug string) (*models.Product, error) {
	product, err := s.productRepo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get product: %w", err)
	}
	return product, nil
}

// ListProducts retrieves products with pagination.
func (s *ProductService) ListProducts(ctx context.Context, pagination utils.PaginationInput) ([]models.Product, int64, error) {
	pagination.Validate()
	return s.productRepo.List(ctx, pagination)
}

// CreateProduct creates a new product with variants.
func (s *ProductService) CreateProduct(ctx context.Context, input CreateProductInput) (*models.Product, error) {
	merchantID := middleware.GetMerchantID(ctx)

	product := &models.Product{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchantID,
		},
		Name:        input.Name,
		Slug:        utils.Slugify(input.Name),
		Description: input.Description,
		BasePrice:   input.BasePrice,
		BrandID:     input.BrandID,
		CategoryID:  input.CategoryID,
		IsPublished: false,
		IsActive:    true,
	}

	if input.ComparePrice != nil {
		product.ComparePrice = input.ComparePrice
	}
	if input.CostPrice != nil {
		product.CostPrice = input.CostPrice
	}

	// Create variants
	for _, v := range input.Variants {
		variant := models.ProductVariant{
			BaseModel: models.BaseModel{ID: uuid.New()},
			ProductID: product.ID,
			SizeID:    v.SizeID,
			ColorID:   v.ColorID,
			SKU:       v.SKU,
			Barcode:   v.Barcode,
			IsActive:  true,
		}
		if v.PriceOverride != nil {
			variant.PriceOverride = v.PriceOverride
		}
		product.Variants = append(product.Variants, variant)
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("create product: %w", err)
	}

	return product, nil
}

// GetProductByBarcode retrieves a variant by barcode (for POS scanning).
func (s *ProductService) GetProductByBarcode(ctx context.Context, barcode string) (*models.ProductVariant, error) {
	return s.productRepo.FindVariantByBarcode(ctx, barcode)
}

// ListSizeGroups retrieves all size groups.
func (s *ProductService) ListSizeGroups(ctx context.Context) ([]models.SizeGroup, error) {
	return s.productRepo.ListSizeGroups(ctx)
}

// CreateSizeGroup creates a new size group with sizes.
func (s *ProductService) CreateSizeGroup(ctx context.Context, name string, sizeNames []string) (*models.SizeGroup, error) {
	merchantID := middleware.GetMerchantID(ctx)

	group := &models.SizeGroup{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchantID,
		},
		Name: name,
	}

	for i, sizeName := range sizeNames {
		group.Sizes = append(group.Sizes, models.Size{
			TenantModel: models.TenantModel{
				ID:         uuid.New(),
				MerchantID: merchantID,
			},
			SizeGroupID: group.ID,
			Name:        sizeName,
			SortOrder:   i,
		})
	}

	if err := s.productRepo.CreateSizeGroup(ctx, group); err != nil {
		return nil, fmt.Errorf("create size group: %w", err)
	}

	return group, nil
}

// CreateProductInput holds input data for creating a product.
type CreateProductInput struct {
	Name         string
	Description  string
	BasePrice    float64
	ComparePrice *float64
	CostPrice    *float64
	BrandID      *uuid.UUID
	CategoryID   *uuid.UUID
	Tags         []string
	Variants     []CreateVariantInput
}

// CreateVariantInput holds input data for creating a variant.
type CreateVariantInput struct {
	SizeID        uuid.UUID
	ColorID       uuid.UUID
	SKU           string
	Barcode       string
	PriceOverride *float64
}
