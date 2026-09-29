package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/naseej/server/internal/models"
)

type contextKey string

const merchantIDKey contextKey = "merchant_id"

// TenantMiddleware extracts X-Merchant-ID from the request header,
// validates the merchant exists and is active, and injects the ID into context.
func TenantMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		merchantIDStr := c.GetHeader("X-Merchant-ID")
		if merchantIDStr == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "X-Merchant-ID header is required",
			})
			return
		}

		merchantID, err := uuid.Parse(merchantIDStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "invalid X-Merchant-ID format",
			})
			return
		}

		// Verify merchant exists and is active
		var merchant models.Merchant
		if err := db.Where("id = ? AND is_active = true", merchantID).First(&merchant).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "merchant not found or inactive",
			})
			return
		}

		// Inject merchant ID into context
		ctx := context.WithValue(c.Request.Context(), merchantIDKey, merchantID)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// GetMerchantID retrieves the merchant ID from the context.
func GetMerchantID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(merchantIDKey).(uuid.UUID); ok {
		return id
	}
	return uuid.Nil
}
