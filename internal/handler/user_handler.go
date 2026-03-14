package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"Inventory-pro/pkg/pagination"
	"github.com/gin-gonic/gin"
	"strconv"
)

type UserHandler struct {
	userHandler service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userHandler: userService}
}
func getIDParam(c *gin.Context, paramName string) (uint, error) {
	idStr := c.Param(paramName)
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GetAllUsers GET /api/admin/users
func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	//check request
	var params request.UserSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	users, total, err := uh.userHandler.SearchAndFilter(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	paginatedResponse := pagination.NewResponse(users, params.Page, params.Limit, total)
	response.Success(c, paginatedResponse)
}

// GetUserById GET /api/admin/users/:id
func (uh *UserHandler) GetUserById(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	user, err := uh.userHandler.GetUserById(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, user)
}

// UpdateUser PUT /api/admin/users/:id
func (uh *UserHandler) UpdateUser(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := uh.userHandler.UpdateUser(id, req); err != nil {
		response.BadRequest(c, err.Error())
	}
	response.Message(c, "User updated")
}

// DeactivateUser PATCH /api/admin/users/:id/deactivate
func (uh *UserHandler) DeactivateUser(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := uh.userHandler.DeactivateUser(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Message(c, "User deactivated")
}

// ActivateUser PATCH /api/admin/users/:id/activate
func (uh *UserHandler) ActivateUser(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := uh.userHandler.ActivateUser(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Message(c, "User activated")
}

// DeleteUser DELETE /api/admin/users/:id
func (uh *UserHandler) DeleteUser(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := uh.userHandler.DeleteUser(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Message(c, "User deleted")
}

// HardDeleteUser DELETE /api/superadmin/users/:id/hard
func (uh *UserHandler) HardDeleteUser(c *gin.Context) {
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := uh.userHandler.HardDeleteUser(id); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Message(c, "User deleted")
}
