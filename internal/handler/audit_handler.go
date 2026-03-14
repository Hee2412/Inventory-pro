package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"Inventory-pro/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type AuditSessionHandler struct {
	service service.AuditSessionService
}

func NewAuditSessionHandler(auditSessionService service.AuditSessionService) *AuditSessionHandler {
	return &AuditSessionHandler{service: auditSessionService}
}

// CreateAuditSession POST api/admin/audit-sessions
func (a *AuditSessionHandler) CreateAuditSession(c *gin.Context) {
	//var request
	var req request.CreateAuditSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//get userID from JWT
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "Invalid token")
		return
	}
	createdBy := userID.(uint)
	//call service
	result, err := a.service.CreateAuditSession(req, createdBy)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// GetAllAuditSession GET /api/superadmin/audit-sessions
func (a *AuditSessionHandler) GetAllAuditSession(c *gin.Context) {
	var params request.AuditSessionSearchParams
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
	sessions, total, err := a.service.GetAllSessionsPaginated(params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	paginatedResponse := pagination.NewResponse(sessions, params.Page, params.Limit, total)
	response.Success(c, paginatedResponse)
}

// GetAuditSessionByID GET /api/superadmin/audit-sessions/:sessionId
func (a *AuditSessionHandler) GetAuditSessionByID(c *gin.Context) {
	id, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := a.service.GetAuditSessionByID(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// AddProductToAudit POST /api/superadmin/audit-sessions/products
func (a *AuditSessionHandler) AddProductToAudit(c *gin.Context) {
	//check request
	var req request.AddProductToAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result, err := a.service.AddProductToAudit(req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, result)
}

// RemoveProductFromAudit DELETE /api/superadmin/audit-sessions/:sessionId/products/:productId
func (a *AuditSessionHandler) RemoveProductFromAudit(c *gin.Context) {
	//get sessionID/productID from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	productID, err := getIDParam(c, "productId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err = a.service.RemoveProductFromAudit(sessionID, productID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Successfully removed product")
}

// CloseAuditSession PATCH /api/superadmin/audit-sessions/:id/close
func (a *AuditSessionHandler) CloseAuditSession(c *gin.Context) {
	//get id from url
	id, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err = a.service.CloseAuditSession(id)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Session has been closed")
}

// UpdateAuditSession PUT /api/superadmin/audit-sessions/:id
func (a *AuditSessionHandler) UpdateAuditSession(c *gin.Context) {
	//get id from URL
	id, err := getIDParam(c, "id")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//bind request
	var req request.UpdateAuditSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	err = a.service.UpdateAuditSession(id, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Updated audit session")
}
