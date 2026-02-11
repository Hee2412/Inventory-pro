package middleware

import (
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
		c.JSON(403, gin.H{"error": "Permission denied"})
		c.Abort()
	}
}
