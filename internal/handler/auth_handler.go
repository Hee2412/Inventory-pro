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
		response.HandleError(c, err)
		return
	}

	token, user, err := handler.authService.Login(req.Username, req.Password)
	if err != nil {
		response.HandleError(c, err)
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
		response.HandleError(c, err)
		return
	}
	creatorID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	creatorRol, _ := c.Get("role")
	roleStr, _ := creatorRol.(string)
	err = handler.authService.Register(&creatorID, roleStr, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, "Register successfully")
}

// GetProfile GET /api/me
func (handler *AuthHandler) GetProfile(c *gin.Context) {
	userID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	user, err := handler.authService.GetProfile(userID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, user)
}
