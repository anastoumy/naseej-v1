package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/naseej/server/internal/models"
)

// StaffRepository handles user/staff and RBAC data access.
type StaffRepository struct {
	*BaseRepository
}

// NewStaffRepository creates a new StaffRepository.
func NewStaffRepository(db *gorm.DB) *StaffRepository {
	return &StaffRepository{BaseRepository: NewBaseRepository(db)}
}

// FindUserByEmail retrieves a user by email within a merchant scope.
func (r *StaffRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	if err := r.Scoped(ctx).
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Locations").
		Where("email = ?", email).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

// FindUserByID retrieves a user by ID within a merchant scope.
func (r *StaffRepository) FindUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.Scoped(ctx).
		Preload("Roles").
		Preload("Roles.Permissions").
		Preload("Locations").
		Where("id = ?", id).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

// FindUserByPin retrieves a user by PIN code within a merchant scope.
func (r *StaffRepository) FindUserByPin(ctx context.Context, pinHash string) (*models.User, error) {
	var user models.User
	if err := r.Scoped(ctx).
		Preload("Roles").
		Preload("Roles.Permissions").
		Where("pin_code = ?", pinHash).
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by pin: %w", err)
	}
	return &user, nil
}

// CreateUser creates a new user.
func (r *StaffRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.DB.Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// ListStaff retrieves all active staff for the merchant.
func (r *StaffRepository) ListStaff(ctx context.Context) ([]models.User, error) {
	var users []models.User
	if err := r.Scoped(ctx).
		Preload("Roles").
		Preload("Locations").
		Where("is_active = true").
		Find(&users).Error; err != nil {
		return nil, fmt.Errorf("list staff: %w", err)
	}
	return users, nil
}

// ListRoles retrieves all roles for the merchant.
func (r *StaffRepository) ListRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	if err := r.Scoped(ctx).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("list roles: %w", err)
	}
	return roles, nil
}

// CreateRole creates a new role.
func (r *StaffRepository) CreateRole(ctx context.Context, role *models.Role) error {
	if err := r.DB.Create(role).Error; err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

// AssignRole assigns a role to a user.
func (r *StaffRepository) AssignRole(ctx context.Context, userID, roleID uuid.UUID) error {
	user, err := r.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	var role models.Role
	if err := r.Scoped(ctx).Where("id = ?", roleID).First(&role).Error; err != nil {
		return fmt.Errorf("find role: %w", err)
	}

	if err := r.DB.Model(user).Association("Roles").Append(&role); err != nil {
		return fmt.Errorf("assign role: %w", err)
	}
	return nil
}
