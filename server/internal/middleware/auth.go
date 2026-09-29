package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const userIDKey contextKey = "user_id"
const userPermissionsKey contextKey = "user_permissions"

// Claims represents the JWT payload.
type Claims struct {
	UserID      uuid.UUID `json:"user_id"`
	MerchantID  uuid.UUID `json:"merchant_id"`
	Permissions []string  `json:"permissions"`
	jwt.RegisteredClaims
}

// AuthMiddleware validates JWT tokens and injects user context.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header required",
			})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}

		// Inject user info into context
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, userIDKey, claims.UserID)
		ctx = context.WithValue(ctx, userPermissionsKey, claims.Permissions)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetUserID retrieves the authenticated user's ID from context.
func GetUserID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(userIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}

// GetUserPermissions retrieves the authenticated user's permissions from context.
func GetUserPermissions(ctx context.Context) []string {
	if perms, ok := ctx.Value(userPermissionsKey).([]string); ok {
		return perms
	}
	return nil
}
