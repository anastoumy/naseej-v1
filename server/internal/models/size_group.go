package models

import "github.com/google/uuid"

// SizeGroup represents a named collection of sizes applicable to a product category.
type SizeGroup struct {
	TenantModel
	Name  string `gorm:"size:100;not null"`
	Sizes []Size `gorm:"foreignKey:SizeGroupID"`
}

// Size represents an individual size within a SizeGroup.
type Size struct {
	TenantModel
	SizeGroupID uuid.UUID `gorm:"type:uuid;not null;index"`
	SizeGroup   SizeGroup `gorm:"foreignKey:SizeGroupID"`
	Name        string    `gorm:"size:20;not null"`
	SortOrder   int       `gorm:"default:0"`
}

// TableName overrides the default table name.
func (SizeGroup) TableName() string {
	return "size_groups"
}

// TableName overrides the default table name.
func (Size) TableName() string {
	return "sizes"
}
