package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/naseej/server/internal/models"
)

// InventoryRepository handles inventory data access.
type InventoryRepository struct {
	*BaseRepository
}

// NewInventoryRepository creates a new InventoryRepository.
func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{BaseRepository: NewBaseRepository(db)}
}

// GetLevel retrieves the inventory level for a variant at a location.
func (r *InventoryRepository) GetLevel(ctx context.Context, variantID, locationID uuid.UUID) (*models.InventoryLevel, error) {
	var level models.InventoryLevel
	if err := r.DB.
		Where("variant_id = ? AND location_id = ?", variantID, locationID).
		First(&level).Error; err != nil {
		return nil, fmt.Errorf("get inventory level: %w", err)
	}
	return &level, nil
}

// GetLevelsByVariant retrieves all inventory levels for a variant across locations.
func (r *InventoryRepository) GetLevelsByVariant(ctx context.Context, variantID uuid.UUID) ([]models.InventoryLevel, error) {
	var levels []models.InventoryLevel
	if err := r.DB.
		Preload("Location").
		Where("variant_id = ?", variantID).
		Find(&levels).Error; err != nil {
		return nil, fmt.Errorf("get levels by variant: %w", err)
	}
	return levels, nil
}

// GetLevelsByLocation retrieves all inventory levels at a location.
func (r *InventoryRepository) GetLevelsByLocation(ctx context.Context, locationID uuid.UUID) ([]models.InventoryLevel, error) {
	var levels []models.InventoryLevel
	if err := r.DB.
		Preload("Variant").
		Preload("Variant.Product").
		Preload("Variant.Size").
		Preload("Variant.Color").
		Where("location_id = ?", locationID).
		Find(&levels).Error; err != nil {
		return nil, fmt.Errorf("get levels by location: %w", err)
	}
	return levels, nil
}

// AdjustStock atomically adjusts stock quantity. Returns error if stock would go negative.
func (r *InventoryRepository) AdjustStock(ctx context.Context, variantID, locationID uuid.UUID, delta int) (*models.InventoryLevel, error) {
	var level models.InventoryLevel

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		// Lock the row for update
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("variant_id = ? AND location_id = ?", variantID, locationID).
			First(&level).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new level if it doesn't exist
				level = models.InventoryLevel{
					VariantID:  variantID,
					LocationID: locationID,
					Quantity:   0,
				}
				if delta < 0 {
					return fmt.Errorf("insufficient stock")
				}
				level.Quantity = delta
				return tx.Create(&level).Error
			}
			return fmt.Errorf("lock inventory level: %w", err)
		}

		newQty := level.Quantity + delta
		if newQty < 0 {
			return fmt.Errorf("insufficient stock: current=%d, delta=%d", level.Quantity, delta)
		}

		level.Quantity = newQty
		return tx.Save(&level).Error
	})

	if err != nil {
		return nil, fmt.Errorf("adjust stock: %w", err)
	}
	return &level, nil
}

// RecordMovement records an inventory movement event.
func (r *InventoryRepository) RecordMovement(ctx context.Context, movement *models.InventoryMovement) error {
	if err := r.DB.Create(movement).Error; err != nil {
		return fmt.Errorf("record movement: %w", err)
	}
	return nil
}
