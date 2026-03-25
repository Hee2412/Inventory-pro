package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"Inventory-pro/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type OrderSessionHandler struct {
	service service.OrderSessionService
}

func NewOrderSessionHandler(orderSessionHandler service.OrderSessionService) *OrderSessionHandler {
	return &OrderSessionHandler{service: orderSessionHandler}
}

// CreateSession POST /api/admin/sessions
func (h *OrderSessionHandler) CreateSession(c *gin.Context) {
	var req request.CreateOrderSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	createdBy, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := h.service.CreateSession(req, createdBy)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, result)
}

// GetAllSessions GET /api/admin/sessions
func (h *OrderSessionHandler) GetAllSessions(c *gin.Context) {
	var params request.SessionSearchParams
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
	sessions, total, err := h.service.GetAllSessions(params)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	paginatedResponse := pagination.NewResponse(sessions, params.Page, params.Limit, total)
	response.Success(c, paginatedResponse)
}

// GetSessionById GET /api/admin/sessions/:sessionId
func (h *OrderSessionHandler) GetSessionById(c *gin.Context) {
	id, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	session, err := h.service.GetSessionById(id)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, session)
}

// AddProductToSession POST /api/admin/sessions/products
func (h *OrderSessionHandler) AddProductToSession(c *gin.Context) {
	var req request.AddProductToSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := h.service.AddProductToSession(req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, result)
}

// RemoveProductFromSession DELETE /api/admin/sessions/:sessionId/products/:productId
func (h *OrderSessionHandler) RemoveProductFromSession(c *gin.Context) {
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	productID, err := getIDParam(c, "productId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	err = h.service.RemoveProductFromSession(sessionID, productID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Product removed from session")
}

// CloseSession PATCH /api/admin/sessions/:sessionId/close
func (h *OrderSessionHandler) CloseSession(c *gin.Context) {
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	err = h.service.CloseSession(sessionID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, "Session closed")
}
