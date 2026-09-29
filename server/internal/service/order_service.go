package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/naseej/server/internal/middleware"
	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/repository"
	"github.com/naseej/server/internal/websocket"
	"github.com/naseej/server/pkg/utils"
)

// OrderService handles order business logic.
type OrderService struct {
	orderRepo     *repository.OrderRepository
	inventoryRepo *repository.InventoryRepository
	wsHub         *websocket.Hub
}

// NewOrderService creates a new OrderService.
func NewOrderService(orderRepo *repository.OrderRepository, inventoryRepo *repository.InventoryRepository, wsHub *websocket.Hub) *OrderService {
	return &OrderService{
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		wsHub:         wsHub,
	}
}

// GetOrder retrieves an order by ID.
func (s *OrderService) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	return s.orderRepo.FindByID(ctx, id)
}

// ListOrders retrieves orders with pagination.
func (s *OrderService) ListOrders(ctx context.Context, pagination utils.PaginationInput) ([]models.Order, int64, error) {
	pagination.Validate()
	return s.orderRepo.List(ctx, pagination)
}

// CreateOrder creates a new order, decrements inventory, and broadcasts.
func (s *OrderService) CreateOrder(ctx context.Context, input CreateOrderInput) (*models.Order, error) {
	merchantID := middleware.GetMerchantID(ctx)
	userID := middleware.GetUserID(ctx)

	order := &models.Order{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchantID,
		},
		OrderNumber: generateOrderNumber(),
		CustomerID:  input.CustomerID,
		LocationID:  input.LocationID,
		StaffID:     &userID,
		Channel:     input.Channel,
		Status:      models.StatusConfirmed,
		Subtotal:    0,
		TaxAmount:   0,
		Total:       0,
	}

	// Process line items
	var subtotal float64
	for _, item := range input.Items {
		lineTotal := item.UnitPrice * float64(item.Quantity)
		orderItem := models.OrderItem{
			BaseModel:    models.BaseModel{ID: uuid.New()},
			OrderID:      order.ID,
			VariantID:    item.VariantID,
			ProductName:  item.ProductName,
			VariantLabel: item.VariantLabel,
			SKU:          item.SKU,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			Total:        lineTotal,
		}
		order.Items = append(order.Items, orderItem)
		subtotal += lineTotal

		// Decrement inventory
		if input.LocationID != nil {
			_, err := s.inventoryRepo.AdjustStock(ctx, item.VariantID, *input.LocationID, -item.Quantity)
			if err != nil {
				return nil, fmt.Errorf("insufficient stock for variant %s: %w", item.VariantID, err)
			}

			// Record inventory movement
			movement := &models.InventoryMovement{
				TenantModel: models.TenantModel{
					ID:         uuid.New(),
					MerchantID: merchantID,
				},
				VariantID:      item.VariantID,
				LocationID:     *input.LocationID,
				QuantityChange: -item.Quantity,
				Reason:         models.MovementSale,
				ReferenceID:    &order.ID,
				ReferenceType:  "order",
				PerformedBy:    &userID,
			}
			if err := s.inventoryRepo.RecordMovement(ctx, movement); err != nil {
				return nil, fmt.Errorf("record sale movement: %w", err)
			}
		}
	}

	order.Subtotal = subtotal
	order.TaxAmount = subtotal * input.TaxRate
	order.Total = subtotal + order.TaxAmount - order.DiscountAmount

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}

	// Broadcast order placed
	if s.wsHub != nil {
		s.wsHub.Broadcast(merchantID, websocket.Event{
			Type:    websocket.EventOrderPlaced,
			Payload: map[string]interface{}{"order_id": order.ID, "total": order.Total},
		})
	}

	return order, nil
}

// UpdateOrderStatus transitions an order to a new status.
func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uuid.UUID, newStatus models.OrderStatus) (*models.Order, error) {
	if err := s.orderRepo.UpdateStatus(ctx, orderID, newStatus); err != nil {
		return nil, fmt.Errorf("update order status: %w", err)
	}
	return s.orderRepo.FindByID(ctx, orderID)
}

func generateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%04d", now.Format("060102"), now.UnixMilli()%10000)
}

// CreateOrderInput holds input for creating an order.
type CreateOrderInput struct {
	CustomerID *uuid.UUID
	LocationID *uuid.UUID
	Channel    models.OrderChannel
	TaxRate    float64
	Items      []CreateOrderItemInput
}

// CreateOrderItemInput holds input for a single order line item.
type CreateOrderItemInput struct {
	VariantID    uuid.UUID
	ProductName  string
	VariantLabel string
	SKU          string
	Quantity     int
	UnitPrice    float64
}
