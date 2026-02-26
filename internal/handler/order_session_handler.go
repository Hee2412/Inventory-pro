package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OrderSessionHandler struct {
	service service.OrderSessionService
}

func NewOrderSessionHandler(orderSessionHandler service.OrderSessionService) *OrderSessionHandler {
	return &OrderSessionHandler{service: orderSessionHandler}
}

// CreateSession POST /api/admin/sessions
func (h *OrderSessionHandler) CreateSession(c *gin.Context) {
	//bind request
	var req request.CreateOrderSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//pull out createdBy from JWT
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	createdBy := userID.(uint)
	//call service
	result, err := h.service.CreateSession(req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//response
	c.JSON(http.StatusCreated, gin.H{
		"message": "session created",
		"data":    result,
	})
}

// GetAllSessions GET /api/admin/sessions
func (h *OrderSessionHandler) GetAllSessions(c *gin.Context) {
	sessions, err := h.service.GetAllSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sessions})
}

// GetSessionById GET /api/admin/sessions/:sessionId
func (h *OrderSessionHandler) GetSessionById(c *gin.Context) {
	//get id param
	id, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	session, err := h.service.GetSessionById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//response
	c.JSON(http.StatusOK, gin.H{"data": session})
}

// AddProductToSession POST /api/admin/sessions/products
func (h *OrderSessionHandler) AddProductToSession(c *gin.Context) {
	//bind request
	var req request.AddProductToSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	err := h.service.AddProductToSession(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//response
	c.JSON(http.StatusCreated, gin.H{"message": "product added to session"})
}

// RemoveProductFromSession DELETE /api/admin/sessions/:sessionId/products/:productId
func (h *OrderSessionHandler) RemoveProductFromSession(c *gin.Context) {
	//get sessionId from url
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get productId from url
	productID, err := getIDParam(c, "productId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	err = h.service.RemoveProductFromSession(sessionID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//response
	c.JSON(http.StatusOK, gin.H{"message": "product removed from session"})
}

// CloseSession PATCH /api/admin/sessions/:sessionId/close
func (h *OrderSessionHandler) CloseSession(c *gin.Context) {
	//get sessionID from url
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	err = h.service.CloseSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "session closed"})
}
