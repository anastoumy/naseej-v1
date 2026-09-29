package repository

import (
	"context"

	"github.com/naseej/server/internal/middleware"
	"gorm.io/gorm"
)

// BaseRepository provides tenant-scoped database access.
type BaseRepository struct {
	DB *gorm.DB
}

// NewBaseRepository creates a new BaseRepository.
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{DB: db}
}

// Scoped returns a GORM session pre-filtered to the current tenant.
// This method MUST be used for every query to prevent cross-tenant data leaks.
func (r *BaseRepository) Scoped(ctx context.Context) *gorm.DB {
	merchantID := middleware.GetMerchantID(ctx)
	return r.DB.Where("merchant_id = ?", merchantID)
}

// Unscoped returns the raw DB session (for non-tenant queries like Merchant lookup).
func (r *BaseRepository) Unscoped() *gorm.DB {
	return r.DB
}
