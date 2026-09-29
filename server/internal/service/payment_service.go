package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/naseej/server/internal/middleware"
	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/repository"
)

// PaymentService handles payment business logic.
type PaymentService struct {
	paymentRepo *repository.PaymentRepository
	orderRepo   *repository.OrderRepository
}

// NewPaymentService creates a new PaymentService.
func NewPaymentService(paymentRepo *repository.PaymentRepository, orderRepo *repository.OrderRepository) *PaymentService {
	return &PaymentService{
		paymentRepo: paymentRepo,
		orderRepo:   orderRepo,
	}
}

// ProcessPayment records a payment for an order.
func (s *PaymentService) ProcessPayment(ctx context.Context, input ProcessPaymentInput) (*models.Payment, error) {
	merchantID := middleware.GetMerchantID(ctx)
	now := time.Now()

	payment := &models.Payment{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchantID,
		},
		OrderID:     input.OrderID,
		Method:      input.Method,
		Status:      models.PayStatusCompleted,
		Amount:      input.Amount,
		Reference:   input.Reference,
		ProcessedAt: &now,
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("process payment: %w", err)
	}

	return payment, nil
}

// RefundPayment issues a refund against a payment.
func (s *PaymentService) RefundPayment(ctx context.Context, input RefundInput) (*models.Refund, error) {
	merchantID := middleware.GetMerchantID(ctx)
	userID := middleware.GetUserID(ctx)

	refund := &models.Refund{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchantID,
		},
		PaymentID:   input.PaymentID,
		OrderID:     input.OrderID,
		Amount:      input.Amount,
		Reason:      input.Reason,
		ProcessedBy: &userID,
	}

	if err := s.paymentRepo.CreateRefund(ctx, refund); err != nil {
		return nil, fmt.Errorf("create refund: %w", err)
	}

	return refund, nil
}

// ProcessPaymentInput holds input for processing a payment.
type ProcessPaymentInput struct {
	OrderID   uuid.UUID
	Method    models.PaymentMethod
	Amount    float64
	Reference string
}

// RefundInput holds input for issuing a refund.
type RefundInput struct {
	PaymentID uuid.UUID
	OrderID   uuid.UUID
	Amount    float64
	Reason    string
}
