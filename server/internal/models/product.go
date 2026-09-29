package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Product represents a base product (parent of variants).
type Product struct {
	TenantModel
	BrandID      *uuid.UUID `gorm:"type:uuid"`
	Brand        *Brand     `gorm:"foreignKey:BrandID"`
	CategoryID   *uuid.UUID `gorm:"type:uuid"`
	Category     *Category  `gorm:"foreignKey:CategoryID"`
	Name         string     `gorm:"size:255;not null"`
	Slug         string     `gorm:"size:255;not null"`
	Description  string
	BasePrice    float64          `gorm:"type:decimal(10,2);not null"`
	ComparePrice *float64         `gorm:"type:decimal(10,2)"`
	CostPrice    *float64         `gorm:"type:decimal(10,2)"`
	IsPublished  bool             `gorm:"default:false"`
	IsActive     bool             `gorm:"default:true"`
	Tags         pq.StringArray   `gorm:"type:text[]"`
	Images       []ProductImage   `gorm:"foreignKey:ProductID"`
	Variants     []ProductVariant `gorm:"foreignKey:ProductID"`
}

// ProductImage represents an image associated with a product.
type ProductImage struct {
	BaseModel
	ProductID uuid.UUID `gorm:"type:uuid;not null;index"`
	URL       string    `gorm:"not null"`
	AltText   string    `gorm:"size:255"`
	SortOrder int       `gorm:"default:0"`
	IsPrimary bool      `gorm:"default:false"`
}

// ProductVariant represents a specific sellable unit (Size × Color combination).
type ProductVariant struct {
	BaseModel
	ProductID     uuid.UUID `gorm:"type:uuid;not null;index"`
	Product       Product   `gorm:"foreignKey:ProductID"`
	SizeID        uuid.UUID `gorm:"type:uuid;not null"`
	Size          Size      `gorm:"foreignKey:SizeID"`
	ColorID       uuid.UUID `gorm:"type:uuid;not null"`
	Color         Color     `gorm:"foreignKey:ColorID"`
	SKU           string    `gorm:"size:50;uniqueIndex;not null"`
	Barcode       string    `gorm:"size:50;uniqueIndex"`
	PriceOverride *float64  `gorm:"type:decimal(10,2)"`
	CostOverride  *float64  `gorm:"type:decimal(10,2)"`
	WeightGrams   *int
	IsActive      bool `gorm:"default:true"`
}

// TableName overrides the default table name.
func (Product) TableName() string {
	return "products"
}

// TableName overrides the default table name.
func (ProductImage) TableName() string {
	return "product_images"
}

// TableName overrides the default table name.
func (ProductVariant) TableName() string {
	return "product_variants"
}
