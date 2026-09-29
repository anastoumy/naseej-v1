package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/naseej/server/internal/middleware"
	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/repository"
	"github.com/naseej/server/internal/websocket"
)

// InventoryService handles inventory business logic.
type InventoryService struct {
	inventoryRepo *repository.InventoryRepository
	wsHub         *websocket.Hub
}

// NewInventoryService creates a new InventoryService.
func NewInventoryService(inventoryRepo *repository.InventoryRepository, wsHub *websocket.Hub) *InventoryService {
	return &InventoryService{
		inventoryRepo: inventoryRepo,
		wsHub:         wsHub,
	}
}

// GetLevels retrieves inventory levels filtered by variant or location.
func (s *InventoryService) GetLevels(ctx context.Context, variantID, locationID *uuid.UUID) ([]models.InventoryLevel, error) {
	if variantID != nil {
		return s.inventoryRepo.GetLevelsByVariant(ctx, *variantID)
	}
	if locationID != nil {
		return s.inventoryRepo.GetLevelsByLocation(ctx, *locationID)
	}
	return nil, fmt.Errorf("either variantId or locationId is required")
}

// AdjustStock adjusts stock and records the movement.
func (s *InventoryService) AdjustStock(ctx context.Context, input AdjustStockInput) (*models.InventoryLevel, error) {
	merchantID := middleware.GetMerchantID(ctx)
	userID := middleware.GetUserID(ctx)

	level, err := s.inventoryRepo.AdjustStock(ctx, input.VariantID, input.LocationID, input.Delta)
	if err != nil {
		return nil, fmt.Errorf("adjust stock: %w", err)
	}

	// Record the movement
	movement := &models.InventoryMovement{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchantID,
		},
		VariantID:      input.VariantID,
		LocationID:     input.LocationID,
		QuantityChange: input.Delta,
		Reason:         input.Reason,
		Notes:          input.Notes,
		PerformedBy:    &userID,
	}

	if err := s.inventoryRepo.RecordMovement(ctx, movement); err != nil {
		return nil, fmt.Errorf("record movement: %w", err)
	}

	// Broadcast stock update via WebSocket
	if s.wsHub != nil {
		s.wsHub.Broadcast(merchantID, websocket.Event{
			Type: websocket.EventStockUpdate,
			Payload: map[string]interface{}{
				"variant_id":  input.VariantID,
				"location_id": input.LocationID,
				"new_qty":     level.Quantity,
			},
		})
	}

	return level, nil
}

// TransferStock moves stock between two locations.
func (s *InventoryService) TransferStock(ctx context.Context, input TransferStockInput) error {
	// Decrement source
	_, err := s.AdjustStock(ctx, AdjustStockInput{
		VariantID:  input.VariantID,
		LocationID: input.FromLocationID,
		Delta:      -input.Quantity,
		Reason:     models.MovementTransferOut,
		Notes:      fmt.Sprintf("Transfer to %s", input.ToLocationID),
	})
	if err != nil {
		return fmt.Errorf("transfer stock out: %w", err)
	}

	// Increment destination
	_, err = s.AdjustStock(ctx, AdjustStockInput{
		VariantID:  input.VariantID,
		LocationID: input.ToLocationID,
		Delta:      input.Quantity,
		Reason:     models.MovementTransferIn,
		Notes:      fmt.Sprintf("Transfer from %s", input.FromLocationID),
	})
	if err != nil {
		return fmt.Errorf("transfer stock in: %w", err)
	}

	return nil
}

// AdjustStockInput holds input for a stock adjustment.
type AdjustStockInput struct {
	VariantID  uuid.UUID
	LocationID uuid.UUID
	Delta      int
	Reason     models.MovementReason
	Notes      string
}

// TransferStockInput holds input for a stock transfer.
type TransferStockInput struct {
	VariantID      uuid.UUID
	FromLocationID uuid.UUID
	ToLocationID   uuid.UUID
	Quantity       int
}
