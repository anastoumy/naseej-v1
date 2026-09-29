package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RBACMiddleware checks if the authenticated user has the required permission.
func RBACMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions := GetUserPermissions(c.Request.Context())
		if !hasPermission(permissions, requiredPermission) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "insufficient permissions",
			})
			return
		}
		c.Next()
	}
}

// HasPermission checks if the given permissions slice contains a specific permission.
func HasPermission(permissions []string, permission string) bool {
	return hasPermission(permissions, permission)
}

func hasPermission(permissions []string, required string) bool {
	for _, p := range permissions {
		if p == required {
			return true
		}
	}
	return false
}

// CheckPermission validates a permission from context (used in resolvers).
func CheckPermission(permissions []string, required string) bool {
	return hasPermission(permissions, required)
}
