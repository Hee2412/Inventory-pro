package middleware

import (
	"Inventory-pro/internal/dto/response"
	"github.com/gin-gonic/gin"
)

func (m *AuthMiddleware) RequireRoles(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		response.Unauthorized(c, "Permission denied")
		c.Abort()
	}
}
