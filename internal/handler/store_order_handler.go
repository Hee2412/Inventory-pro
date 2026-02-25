package handler

import (
	"Inventory-pro/internal/dto/request"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get storeId
	userID, exists := c.Get("storeId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	storeID := userID.(uint)
	//call service
	result, err := s.service.GetOrCreateOrder(sessionID, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
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
	result := s.service.UpdateOrder(orderID, req)
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// SubmitOrder POST   /api/store/orders/:orderId/submit
func (s *StoreOrderHandler) SubmitOrder(c *gin.Context) {
	orderID, err := getIDParam(c, "orderId")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, exists := c.Get("store_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	//call service
	err = s.service.SubmitOrder(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": "Submit order successfully"})
}

// GetOrderDetail GET    /api/store/orders/:orderId
func (s *StoreOrderHandler) GetOrderDetail(c *gin.Context) {
	orderID, err := getIDParam(c, "orderId")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	result, err := s.service.GetOrderDetail(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetMyOrder GET    /api/store/orders
func (s *StoreOrderHandler) GetMyOrder(c *gin.Context) {
	userId, exist := c.Get("storeId")
	if !exist {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	storeID := userId.(uint)
	result, err := s.service.GetMyOrder(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
