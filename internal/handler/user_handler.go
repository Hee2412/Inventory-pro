package handler

import (
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	userHandler service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userHandler: userService}
}

func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := uh.userHandler.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fails to get users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (uh *UserHandler) GetUserById(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "userId not exists"})
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
	user, err := uh.userHandler.GetUserById(uint(floatID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unfounded User"})
		return
	}
	c.JSON(http.StatusOK, user)
}
