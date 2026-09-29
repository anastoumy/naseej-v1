package models

import "github.com/google/uuid"

// Category represents a hierarchical product category.
type Category struct {
	TenantModel
	ParentID    *uuid.UUID `gorm:"type:uuid;index"`
	Parent      *Category  `gorm:"foreignKey:ParentID"`
	Children    []Category `gorm:"foreignKey:ParentID"`
	SizeGroupID *uuid.UUID `gorm:"type:uuid"`
	SizeGroup   *SizeGroup `gorm:"foreignKey:SizeGroupID"`
	Name        string     `gorm:"size:100;not null"`
	Slug        string     `gorm:"size:100;not null"`
	Description string
	ImageURL    string
	SortOrder   int `gorm:"default:0"`
}

// Color represents a color option for product variants.
type Color struct {
	TenantModel
	Name    string `gorm:"size:50;not null"`
	HexCode string `gorm:"size:7"`
}

// Brand represents a product brand.
type Brand struct {
	TenantModel
	Name    string `gorm:"size:100;not null"`
	Slug    string `gorm:"size:100;not null"`
	LogoURL string
}

// TableName overrides the default table name.
func (Category) TableName() string {
	return "categories"
}

// TableName overrides the default table name.
func (Color) TableName() string {
	return "colors"
}

// TableName overrides the default table name.
func (Brand) TableName() string {
	return "brands"
}
