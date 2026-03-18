package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
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
		response.HandleError(c, err)
		return
	}
	//get storeId
	userID, exists := c.Get("userId")
	if !exists {
		response.HandleError(c, err)
		return
	}
	storeID := userID.(uint)
	//call service
	result, err := s.service.GetOrCreateOrder(sessionID, storeID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// UpdateOrder PUT    /api/store/orders/:orderId/items
func (s *StoreOrderHandler) UpdateOrder(c *gin.Context) {
	var req request.UpdateOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	//get orderId
	orderID, err := getIDParam(c, "orderId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	//call service
	err = s.service.UpdateOrder(orderID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Order updated")
}

// GetOrderDetail GET    /api/store/orders/:orderId
func (s *StoreOrderHandler) GetOrderDetail(c *gin.Context) {
	orderID, err := getIDParam(c, "orderId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	//call service
	result, err := s.service.GetOrderDetail(orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetMyOrder GET    /api/store/orders
func (s *StoreOrderHandler) GetMyOrder(c *gin.Context) {
	storeID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := s.service.GetMyOrder(storeID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// UpdateStatus PUT /api/store/orders/:orderId
func (s *StoreOrderHandler) UpdateStatus(c *gin.Context) {
	orderID, err := getIDParam(c, "orderId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	updatedOrder, err := s.service.UpdateStatus(orderID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	msg := "Status updated to DRAFT"
	if updatedOrder.Status == "NO_ORDER" {
		msg = "Confirmed no order"
	}
	response.Message(c, msg)
}
