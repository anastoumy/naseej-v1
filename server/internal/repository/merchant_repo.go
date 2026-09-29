package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/naseej/server/internal/models"
)

// MerchantRepository handles merchant data access.
type MerchantRepository struct {
	*BaseRepository
}

// NewMerchantRepository creates a new MerchantRepository.
func NewMerchantRepository(db *gorm.DB) *MerchantRepository {
	return &MerchantRepository{BaseRepository: NewBaseRepository(db)}
}

// FindByID retrieves a merchant by ID.
func (r *MerchantRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.Unscoped().Where("id = ?", id).First(&merchant).Error; err != nil {
		return nil, fmt.Errorf("find merchant by id: %w", err)
	}
	return &merchant, nil
}

// FindBySlug retrieves a merchant by slug.
func (r *MerchantRepository) FindBySlug(ctx context.Context, slug string) (*models.Merchant, error) {
	var merchant models.Merchant
	if err := r.Unscoped().Where("slug = ?", slug).First(&merchant).Error; err != nil {
		return nil, fmt.Errorf("find merchant by slug: %w", err)
	}
	return &merchant, nil
}

// GetLocations retrieves all active locations for a merchant.
func (r *MerchantRepository) GetLocations(ctx context.Context) ([]models.Location, error) {
	var locations []models.Location
	if err := r.Scoped(ctx).Where("is_active = true").Find(&locations).Error; err != nil {
		return nil, fmt.Errorf("get locations: %w", err)
	}
	return locations, nil
}
