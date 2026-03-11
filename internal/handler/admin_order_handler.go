package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"Inventory-pro/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type AdminOrderHandler struct {
	service service.AdminOrderService
}

func NewAdminOrderHandler(adminOrderHandler service.AdminOrderService) *AdminOrderHandler {
	return &AdminOrderHandler{service: adminOrderHandler}
}

// GetAllOrderInSession GET    /api/admin/sessions/:sessionId/orders
func (a *AdminOrderHandler) GetAllOrderInSession(c *gin.Context) {
	//get sessionID
	sessionId, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	orders, err := a.service.GetAllOrderInSession(sessionId)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, orders)
}

// ApproveOrder POST   /api/admin/orders/:orderId/approve
func (a *AdminOrderHandler) ApproveOrder(c *gin.Context) {
	orderId, err := getIDParam(c, "orderId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err = a.service.ApproveOrder(orderId)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Order approved successfully")
}

// DeclineOrder POST   /api/admin/orders/:orderId/decline
func (a *AdminOrderHandler) DeclineOrder(c *gin.Context) {
	//get orderID
	orderId, err := getIDParam(c, "orderId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//var request pull out reason
	var req request.DeclineOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	err = a.service.DeclineOrder(orderId, req.Reason)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Order declined successfully")
}

func (a *AdminOrderHandler) GetAllOrders(c *gin.Context) {
	param := pagination.ParseParams(c)
	//check request
	var params request.OrderSearchParams
	if err := c.ShouldBindQuery(&param); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	orders, total, err := a.service.GetAllPaginatedSessions(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	paginatedResponse := pagination.NewResponse(orders, params.Page, params.Limit, total)
	response.Success(c, paginatedResponse)
}
