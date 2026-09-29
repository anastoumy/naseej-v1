package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/pkg/utils"
)

// CustomerRepository handles customer data access.
type CustomerRepository struct {
	*BaseRepository
}

// NewCustomerRepository creates a new CustomerRepository.
func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{BaseRepository: NewBaseRepository(db)}
}

// FindByID retrieves a customer by ID.
func (r *CustomerRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := r.Scoped(ctx).
		Preload("Addresses").
		Where("id = ?", id).
		First(&customer).Error; err != nil {
		return nil, fmt.Errorf("find customer by id: %w", err)
	}
	return &customer, nil
}

// Search searches customers by name, email, or phone.
func (r *CustomerRepository) Search(ctx context.Context, query string, pagination utils.PaginationInput) ([]models.Customer, int64, error) {
	var customers []models.Customer
	var total int64

	dbQuery := r.Scoped(ctx).Model(&models.Customer{})

	if query != "" {
		search := "%" + query + "%"
		dbQuery = dbQuery.Where(
			"first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
			search, search, search, search,
		)
	}

	dbQuery.Count(&total)

	if err := dbQuery.
		Preload("Addresses").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order("created_at DESC").
		Find(&customers).Error; err != nil {
		return nil, 0, fmt.Errorf("search customers: %w", err)
	}
	return customers, total, nil
}

// Create creates a new customer.
func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	if err := r.DB.Create(customer).Error; err != nil {
		return fmt.Errorf("create customer: %w", err)
	}
	return nil
}

// Update updates an existing customer.
func (r *CustomerRepository) Update(ctx context.Context, customer *models.Customer) error {
	if err := r.Scoped(ctx).Save(customer).Error; err != nil {
		return fmt.Errorf("update customer: %w", err)
	}
	return nil
}
