package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authServer service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authServer}
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	token, user, err := handler.authService.Login(req.Username, req.Password)
	if err != nil {
		response.Unauthorized(c, "Invalid token")
		return
	}

	result := response.LoginResponse{
		Token:     token,
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		StoreName: user.StoreName,
		StoreCode: user.StoreCode,
	}
	response.Success(c, result)
}

// Register POST /api/admin/register
func (handler *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	id, exists := c.Get("userId")
	if !exists {
		response.BadRequest(c, "user id not found")
		return
	}
	creatorID := id.(uint)
	creatorRol, _ := c.Get("role")
	roleStr, _ := creatorRol.(string)
	err := handler.authService.Register(&creatorID, roleStr, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, "Register successfully")
}

// GetProfile GET /api/me
func (handler *AuthHandler) GetProfile(c *gin.Context) {
	id, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "Invalid token")
		return
	}
	userId := id.(uint)
	user, err := handler.authService.GetProfile(userId)
	if err != nil {
		response.NotFound(c, "User not found")
		return
	}
	response.Success(c, user)
}
