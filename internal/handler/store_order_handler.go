package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StoreOrderHandler struct {
	service service.StoreOrderService
}

func NewStoreOrderHandler(storeOrderHandler service.StoreOrderService) *StoreOrderHandler {
	return &StoreOrderHandler{service: storeOrderHandler}
}

// GetOrCreateOrder GET    /api/store/sessions/:sessionId/order
func (s *StoreOrderHandler) GetOrCreateOrder(c *gin.Context) {
	//get sessionID
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//get storeId
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "Invalid token")
		return
	}
	storeID := userID.(uint)
	//call service
	result, err := s.service.GetOrCreateOrder(sessionID, storeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// UpdateOrder PUT    /api/store/orders/:orderId/items
func (s *StoreOrderHandler) UpdateOrder(c *gin.Context) {
	var req request.UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get orderId
	orderID, err := getIDParam(c, "orderId")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	err = s.service.UpdateOrder(orderID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Order updated")
}

// GetOrderDetail GET    /api/store/orders/:orderId
func (s *StoreOrderHandler) GetOrderDetail(c *gin.Context) {
	orderID, err := getIDParam(c, "orderId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	result, err := s.service.GetOrderDetail(orderID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetMyOrder GET    /api/store/orders
func (s *StoreOrderHandler) GetMyOrder(c *gin.Context) {
	userId, exist := c.Get("userId")
	if !exist {
		response.Unauthorized(c, "Invalid token")
		return
	}
	storeID := userId.(uint)
	result, err := s.service.GetMyOrder(storeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, result)
}
