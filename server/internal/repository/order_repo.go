package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/pkg/utils"
)

// OrderRepository handles order data access.
type OrderRepository struct {
	*BaseRepository
}

// NewOrderRepository creates a new OrderRepository.
func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{BaseRepository: NewBaseRepository(db)}
}

// FindByID retrieves an order by ID.
func (r *OrderRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.Scoped(ctx).
		Preload("Customer").
		Preload("Location").
		Preload("Staff").
		Preload("Items").
		Preload("Items.Variant").
		Preload("Payments").
		Where("id = ?", id).
		First(&order).Error; err != nil {
		return nil, fmt.Errorf("find order by id: %w", err)
	}
	return &order, nil
}

// List retrieves orders with pagination.
func (r *OrderRepository) List(ctx context.Context, pagination utils.PaginationInput) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64

	query := r.Scoped(ctx).Model(&models.Order{})
	query.Count(&total)

	if err := query.
		Preload("Customer").
		Preload("Location").
		Preload("Items").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("list orders: %w", err)
	}
	return orders, total, nil
}

// Create creates a new order with items.
func (r *OrderRepository) Create(ctx context.Context, order *models.Order) error {
	if err := r.DB.Create(order).Error; err != nil {
		return fmt.Errorf("create order: %w", err)
	}
	return nil
}

// UpdateStatus updates the status of an order.
func (r *OrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.OrderStatus) error {
	if err := r.Scoped(ctx).Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("update order status: %w", err)
	}
	return nil
}
