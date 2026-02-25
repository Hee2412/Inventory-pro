package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orders, err := a.service.GetAllOrderInSession(sessionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// ApproveOrder POST   /api/admin/orders/:orderId/approve
func (a *AdminOrderHandler) ApproveOrder(c *gin.Context) {
	orderId, err := getIDParam(c, "orderId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = a.service.ApproveOrder(orderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "Order approved successfully"})
}

// DeclineOrder POST   /api/admin/orders/:orderId/decline
func (a *AdminOrderHandler) DeclineOrder(c *gin.Context) {
	//get orderID
	orderId, err := getIDParam(c, "orderId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//var request pull out reason
	var req request.DeclineOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	err = a.service.DeclineOrder(orderId, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "Order declined successfully"})
}
