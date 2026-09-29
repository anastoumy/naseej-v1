package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/naseej/server/internal/models"
)

// PaymentRepository handles payment data access.
type PaymentRepository struct {
	*BaseRepository
}

// NewPaymentRepository creates a new PaymentRepository.
func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{BaseRepository: NewBaseRepository(db)}
}

// Create creates a new payment record.
func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	if err := r.DB.Create(payment).Error; err != nil {
		return fmt.Errorf("create payment: %w", err)
	}
	return nil
}

// FindByOrderID retrieves all payments for an order.
func (r *PaymentRepository) FindByOrderID(ctx context.Context, orderID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	if err := r.Scoped(ctx).Where("order_id = ?", orderID).Find(&payments).Error; err != nil {
		return nil, fmt.Errorf("find payments by order: %w", err)
	}
	return payments, nil
}

// CreateRefund creates a refund record.
func (r *PaymentRepository) CreateRefund(ctx context.Context, refund *models.Refund) error {
	if err := r.DB.Create(refund).Error; err != nil {
		return fmt.Errorf("create refund: %w", err)
	}
	return nil
}
