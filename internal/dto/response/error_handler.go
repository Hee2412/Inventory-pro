package response

import (
	"Inventory-pro/internal/domain"
	"errors"
	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error) {
	if err == nil {
		return
	}

	switch {
	// 404 Not Found
	case errors.Is(err, domain.ErrUserNotFound):
		NotFound(c, "User not found")
	case errors.Is(err, domain.ErrProductNotFound):
		NotFound(c, "Product not found")
	case errors.Is(err, domain.ErrSessionNotFound):
		NotFound(c, "Session not found")
	case errors.Is(err, domain.ErrOrderNotFound):
		NotFound(c, "Order not found")
	case errors.Is(err, domain.ErrAuditNotFound):
		NotFound(c, "Audit not found")

	// 401 Unauthorized
	case errors.Is(err, domain.ErrInvalidCredentials):
		Unauthorized(c, "Invalid credentials")
	case errors.Is(err, domain.ErrUnauthorized):
		Unauthorized(c, "Unauthorized")

	// 403 Forbidden
	case errors.Is(err, domain.ErrForbidden):
		Error(c, 403, "Forbidden")
	case errors.Is(err, domain.ErrUserInactive):
		Error(c, 403, "User is inactive")

	// 400 Bad Request
	case errors.Is(err, domain.ErrInvalidInput):
		BadRequest(c, err.Error())
	case errors.Is(err, domain.ErrSessionClosed):
		BadRequest(c, "Session is closed")
	case errors.Is(err, domain.ErrOrderClosed):
		BadRequest(c, "Cannot update closed order")
	case errors.Is(err, domain.ErrEmptyOrder):
		BadRequest(c, "Cannot submit empty order")
	case errors.Is(err, domain.ErrAuditExpired):
		BadRequest(c, "Audit session has expired")

	// 500 Internal Server Error
	default:
		InternalError(c, "Internal server error")
	}
}
