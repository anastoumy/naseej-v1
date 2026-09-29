package models

import "github.com/google/uuid"

// InventoryLevel tracks stock quantity per variant per location.
type InventoryLevel struct {
	BaseModel
	VariantID    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Variant      ProductVariant `gorm:"foreignKey:VariantID"`
	LocationID   uuid.UUID      `gorm:"type:uuid;not null;index"`
	Location     Location       `gorm:"foreignKey:LocationID"`
	Quantity     int            `gorm:"not null;default:0;check:quantity >= 0"`
	Reserved     int            `gorm:"not null;default:0;check:reserved >= 0"`
	ReorderPoint int            `gorm:"default:0"`
}

// MovementReason categorizes why stock changed.
type MovementReason string

const (
	MovementSale        MovementReason = "sale"
	MovementReturn      MovementReason = "return"
	MovementAdjustment  MovementReason = "adjustment"
	MovementTransferIn  MovementReason = "transfer_in"
	MovementTransferOut MovementReason = "transfer_out"
	MovementRestock     MovementReason = "restock"
	MovementDamaged     MovementReason = "damaged"
	MovementLost        MovementReason = "lost"
	MovementManualCount MovementReason = "manual_count"
)

// InventoryMovement records a single stock change event.
type InventoryMovement struct {
	TenantModel
	VariantID      uuid.UUID      `gorm:"type:uuid;not null;index"`
	Variant        ProductVariant `gorm:"foreignKey:VariantID"`
	LocationID     uuid.UUID      `gorm:"type:uuid;not null"`
	Location       Location       `gorm:"foreignKey:LocationID"`
	QuantityChange int            `gorm:"not null"`
	Reason         MovementReason `gorm:"size:20;not null"`
	ReferenceID    *uuid.UUID     `gorm:"type:uuid"`
	ReferenceType  string         `gorm:"size:50"`
	Notes          string
	PerformedBy    *uuid.UUID `gorm:"type:uuid"`
	Performer      *User      `gorm:"foreignKey:PerformedBy"`
}

// TableName overrides the default table name.
func (InventoryLevel) TableName() string {
	return "inventory_levels"
}

// TableName overrides the default table name.
func (InventoryMovement) TableName() string {
	return "inventory_movements"
}
