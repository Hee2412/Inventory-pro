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
	//bind request
	var req request.CreateOrderSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//pull out createdBy from JWT
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "Invalid token")
		return
	}
	createdBy := userID.(uint)
	//call service
	result, err := h.service.CreateSession(req, createdBy)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	//response
	response.Created(c, result)
}

// GetAllSessions GET /api/admin/sessions
func (h *OrderSessionHandler) GetAllSessions(c *gin.Context) {
	param := pagination.ParseParams(c)
	//check request
	var params request.SessionSearchParams
	if err := c.ShouldBindQuery(&param); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	sessions, total, err := h.service.GetAllSessions(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	//create response
	paginatedResponse := pagination.NewResponse(sessions, params.Page, params.Limit, total)
	response.Success(c, paginatedResponse)
}

// GetSessionById GET /api/admin/sessions/:sessionId
func (h *OrderSessionHandler) GetSessionById(c *gin.Context) {
	//get id param
	id, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	session, err := h.service.GetSessionById(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	//response
	response.Success(c, session)
}

// AddProductToSession POST /api/admin/sessions/products
func (h *OrderSessionHandler) AddProductToSession(c *gin.Context) {
	//bind request
	var req request.AddProductToSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	result, err := h.service.AddProductToSession(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	//response
	response.Created(c, result)
}

// RemoveProductFromSession DELETE /api/admin/sessions/:sessionId/products/:productId
func (h *OrderSessionHandler) RemoveProductFromSession(c *gin.Context) {
	//get sessionId from url
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//get productId from url
	productID, err := getIDParam(c, "productId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	err = h.service.RemoveProductFromSession(sessionID, productID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	//response
	response.Message(c, "Product removed from session")
}

// CloseSession PATCH /api/admin/sessions/:sessionId/close
func (h *OrderSessionHandler) CloseSession(c *gin.Context) {
	//get sessionID from url
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	err = h.service.CloseSession(sessionID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, "Session closed")
}
