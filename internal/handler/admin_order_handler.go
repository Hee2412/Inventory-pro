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
	sessionId, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	orders, err := a.service.GetAllOrderInSession(sessionId)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, orders)
}

// ApproveOrder POST   /api/admin/orders/:orderId/approve
func (a *AdminOrderHandler) ApproveOrder(c *gin.Context) {
	orderId, err := getIDParam(c, "orderId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	err = a.service.ApproveOrder(orderId)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Order approved successfully")
}

// DeclineOrder POST   /api/admin/orders/:orderId/decline
func (a *AdminOrderHandler) DeclineOrder(c *gin.Context) {
	orderId, err := getIDParam(c, "orderId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	var req request.DeclineOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	err = a.service.DeclineOrder(orderId, req.Reason)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Order declined successfully")
}

func (a *AdminOrderHandler) GetAllOrders(c *gin.Context) {
	var params request.OrderSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.HandleError(c, err)
		return
	}
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Limit <= 0 {
		params.Limit = 20
	}
	orders, total, err := a.service.GetAllPaginatedSessions(params)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	paginatedResponse := pagination.NewResponse(orders, params.Page, params.Limit, total)
	response.Success(c, paginatedResponse)
}

// DeliverOrder POST /api/admin/orders/:orderId/deliver
func (a *AdminOrderHandler) DeliverOrder(c *gin.Context) {
	orderId, err := getIDParam(c, "orderId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	if err = a.service.DeliverOrder(orderId); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Order marked as delivered")
}

// RedeliverOrder POST /api/admin/orders/:orderId/redeliver
func (a *AdminOrderHandler) RedeliverOrder(c *gin.Context) {
	orderId, err := getIDParam(c, "orderId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	var req request.RedeliverRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	if err = a.service.RedeliverOrder(orderId, req); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Order redelivered with updated quantities")
}
