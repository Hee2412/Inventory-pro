package handler

import (
	"Inventory-pro/internal/domain"
	"github.com/gin-gonic/gin"
	"strconv"
)

// getIDParam parses URL parameter (:id, :orderId, etc.) to uint
func getIDParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, domain.ErrInvalidInput
	}
	return uint(id), nil
}

// getUserID gets authenticated user ID from JWT context
func getUserID(c *gin.Context) (uint, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return 0, domain.ErrUnauthorized
	}
	id, ok := userId.(uint)
	if !ok {
		return 0, domain.ErrUnauthorized
	}
	return id, nil
}
