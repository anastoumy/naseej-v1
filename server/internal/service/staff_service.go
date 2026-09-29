package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/naseej/server/internal/middleware"
	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/repository"
)

// StaffService handles staff and RBAC business logic.
type StaffService struct {
	staffRepo *repository.StaffRepository
}

// NewStaffService creates a new StaffService.
func NewStaffService(staffRepo *repository.StaffRepository) *StaffService {
	return &StaffService{staffRepo: staffRepo}
}

// ListStaff retrieves all active staff members.
func (s *StaffService) ListStaff(ctx context.Context) ([]models.User, error) {
	return s.staffRepo.ListStaff(ctx)
}

// ListRoles retrieves all roles.
func (s *StaffService) ListRoles(ctx context.Context) ([]models.Role, error) {
	return s.staffRepo.ListRoles(ctx)
}

// CreateStaff creates a new staff member.
func (s *StaffService) CreateStaff(ctx context.Context, input CreateStaffInput) (*models.User, error) {
	merchantID := middleware.GetMerchantID(ctx)

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		TenantModel: models.TenantModel{
			ID:         uuid.New(),
			MerchantID: merchantID,
		},
		Email:        input.Email,
		PasswordHash: passwordHash,
		FirstName:    input.FirstName,
		LastName:     input.LastName,
		Phone:        input.Phone,
		IsActive:     true,
	}

	if input.Pin != "" {
		pinHash, err := HashPin(input.Pin)
		if err != nil {
			return nil, fmt.Errorf("hash pin: %w", err)
		}
		user.PinCode = pinHash
	}

	if err := s.staffRepo.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("create staff: %w", err)
	}

	// Assign role if provided
	if input.RoleID != nil {
		if err := s.staffRepo.AssignRole(ctx, user.ID, *input.RoleID); err != nil {
			return nil, fmt.Errorf("assign role: %w", err)
		}
	}

	return user, nil
}

// AssignRole assigns a role to a staff member.
func (s *StaffService) AssignRole(ctx context.Context, userID, roleID uuid.UUID) (*models.User, error) {
	if err := s.staffRepo.AssignRole(ctx, userID, roleID); err != nil {
		return nil, fmt.Errorf("assign role: %w", err)
	}
	return s.staffRepo.FindUserByID(ctx, userID)
}

// CreateStaffInput holds input for creating a staff member.
type CreateStaffInput struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	Phone     string
	Pin       string
	RoleID    *uuid.UUID
}
