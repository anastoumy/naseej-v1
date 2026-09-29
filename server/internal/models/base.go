package models

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel is embedded in non-tenant entities (like Merchant itself).
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// TenantModel is embedded in all merchant-scoped entities.
type TenantModel struct {
	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	MerchantID uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
