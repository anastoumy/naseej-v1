package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/naseej/server/internal/middleware"
	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/repository"
	"github.com/naseej/server/pkg/utils"
)

// CustomerService handles customer business logic.
type CustomerService struct {
	customerRepo *repository.CustomerRepository
}

// NewCustomerService creates a new CustomerService.
func NewCustomerService(customerRepo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{customerRepo: customerRepo}
}

// GetCustomer retrieves a customer by ID.
func (s *CustomerService) GetCustomer(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	return s.customerRepo.FindByID(ctx, id)
}

// SearchCustomers searches customers by query string.
func (s *CustomerService) SearchCustomers(ctx context.Context, query string, pagination utils.PaginationInput) ([]models.Customer, int64, error) {
	pagination.Validate()
	return s.customerRepo.Search(ctx, query, pagination)
}

// CreateCustomer creates a new customer.
func (s *CustomerService) CreateCustomer(ctx context.Context, input CreateCustomerInput) (*models.Customer, error) {
	merchantID := middleware.GetMerchantID(ctx)

	customer := &models.Customer{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchantID,
		},
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
		Phone:     input.Phone,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("create customer: %w", err)
	}

	return customer, nil
}

// UpdateCustomer updates an existing customer.
func (s *CustomerService) UpdateCustomer(ctx context.Context, id uuid.UUID, input UpdateCustomerInput) (*models.Customer, error) {
	customer, err := s.customerRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find customer: %w", err)
	}

	if input.FirstName != "" {
		customer.FirstName = input.FirstName
	}
	if input.LastName != "" {
		customer.LastName = input.LastName
	}
	if input.Email != "" {
		customer.Email = input.Email
	}
	if input.Phone != "" {
		customer.Phone = input.Phone
	}

	if err := s.customerRepo.Update(ctx, customer); err != nil {
		return nil, fmt.Errorf("update customer: %w", err)
	}

	return customer, nil
}

// CreateCustomerInput holds input for creating a customer.
type CreateCustomerInput struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}

// UpdateCustomerInput holds input for updating a customer.
type UpdateCustomerInput struct {
	FirstName string
	LastName  string
	Email     string
	Phone     string
}
