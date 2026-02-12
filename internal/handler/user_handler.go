package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserHandler struct {
	userHandler service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userHandler: userService}
}
func getIDParam(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GET /api/admin/users

func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := uh.userHandler.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fails to get users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

// GetUserById GET /api/admin/users/:id
func (uh *UserHandler) GetUserById(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := uh.userHandler.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": user})
}

// UpdateUser PUT /api/admin/users/:id
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if err := uh.userHandler.UpdateUser(id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated user"})
}

// DeactivateUser PATCH /api/admin/users/:id/deactivate
func (uh *UserHandler) DeactivateUser(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := uh.userHandler.DeactivateUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deactivated user"})
}

// ActivateUser PATCH /api/admin/users/:id/activate
func (uh *UserHandler) ActivateUser(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := uh.userHandler.ActivateUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully activated user"})
}

// DeleteUser DELETE /api/admin/users/:id
func (uh *UserHandler) DeleteUser(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := uh.userHandler.DeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted user"})
}

// HardDeleteUser DELETE /api/superadmin/users/:id/hard
func (uh *UserHandler) HardDeleteUser(c *gin.Context) {
	id, err := getIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	if err := uh.userHandler.HardDeleteUser(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted user"})
}
