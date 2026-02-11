package handler

import (
	"Inventory-pro/internal/dto"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authServer service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authServer}
}

func (handler *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	token, user, err := handler.authService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.LoginResponse{
		Token:     token,
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		StoreName: user.StoreName,
		StoreCode: user.StoreCode,
	})

}

func (handler *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	creatorRol, _ := c.Get("role")
	roleStr, _ := creatorRol.(string)
	err := handler.authService.Register(roleStr, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "register successfully",
	})
}

func (handler *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not Found"})
		return
	}
	floatID, ok := userID.(float64)
	if !ok {
		uID, okeUint := userID.(uint)
		if !okeUint {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error format"})
			return
		}
		floatID = float64(uID)
	}
	user, err := handler.authService.GetProfile(uint(floatID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unfounded User"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user,
	})
}
