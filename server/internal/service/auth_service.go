package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/naseej/server/config"
	"github.com/naseej/server/internal/middleware"
	"github.com/naseej/server/internal/models"
	"github.com/naseej/server/internal/repository"
)

// AuthService handles authentication and JWT operations.
type AuthService struct {
	staffRepo *repository.StaffRepository
	config    *config.Config
}

// NewAuthService creates a new AuthService.
func NewAuthService(staffRepo *repository.StaffRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		staffRepo: staffRepo,
		config:    cfg,
	}
}

// AuthPayload is returned on successful authentication.
type AuthPayload struct {
	Token        string
	RefreshToken string
	User         *models.User
}

// Login authenticates a user with email and password.
func (s *AuthService) Login(ctx context.Context, email, password string) (*AuthPayload, error) {
	user, err := s.staffRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !user.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	token, refreshToken, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &AuthPayload{
		Token:        token,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// LoginWithPin authenticates a user with a PIN code (POS quick login).
func (s *AuthService) LoginWithPin(ctx context.Context, pin string) (*AuthPayload, error) {
	// Hash the pin and look up user
	users, err := s.staffRepo.ListStaff(ctx)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	var matchedUser *models.User
	for i := range users {
		if users[i].PinCode != "" {
			if err := bcrypt.CompareHashAndPassword([]byte(users[i].PinCode), []byte(pin)); err == nil {
				matchedUser = &users[i]
				break
			}
		}
	}

	if matchedUser == nil {
		return nil, fmt.Errorf("invalid PIN")
	}

	if !matchedUser.IsActive {
		return nil, fmt.Errorf("account is deactivated")
	}

	token, refreshToken, err := s.generateTokens(ctx, matchedUser)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &AuthPayload{
		Token:        token,
		RefreshToken: refreshToken,
		User:         matchedUser,
	}, nil
}

// RefreshToken validates a refresh token and issues a new token pair.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (*AuthPayload, error) {
	claims := &middleware.Claims{}
	token, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.Auth.JWTSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid refresh token")
	}

	user, err := s.staffRepo.FindUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	newToken, newRefresh, err := s.generateTokens(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return &AuthPayload{
		Token:        newToken,
		RefreshToken: newRefresh,
		User:         user,
	}, nil
}

// HashPassword hashes a plaintext password with bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}
	return string(bytes), nil
}

// HashPin hashes a PIN code with bcrypt.
func HashPin(pin string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash pin: %w", err)
	}
	return string(bytes), nil
}

func (s *AuthService) generateTokens(ctx context.Context, user *models.User) (string, string, error) {
	merchantID := middleware.GetMerchantID(ctx)

	// Collect all permissions from user's roles
	permissions := s.collectPermissions(user)

	// Access token
	accessClaims := &middleware.Claims{
		UserID:      user.ID,
		MerchantID:  merchantID,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.Auth.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return "", "", fmt.Errorf("sign access token: %w", err)
	}

	// Refresh token (longer lived, minimal claims)
	refreshClaims := &middleware.Claims{
		UserID:     user.ID,
		MerchantID: merchantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.Auth.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
	}
	refreshTokenJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshString, err := refreshTokenJWT.SignedString([]byte(s.config.Auth.JWTSecret))
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	return accessString, refreshString, nil
}

func (s *AuthService) collectPermissions(user *models.User) []string {
	permSet := make(map[string]bool)
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			key := perm.Resource + ":" + perm.Action
			permSet[key] = true
		}
	}

	permissions := make([]string, 0, len(permSet))
	for p := range permSet {
		permissions = append(permissions, p)
	}
	return permissions
}
