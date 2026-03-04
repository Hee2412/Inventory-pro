package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get userID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	createdBy := userID.(uint)
	//call service
	result, err := a.service.CreateAuditSession(req, createdBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetAllAuditSession GET /api/superadmin/audit-sessions
func (a *AuditSessionHandler) GetAllAuditSession(c *gin.Context) {
	// cal service
	result, err := a.service.GetAllAuditSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetAuditSessionByID GET /api/superadmin/audit-sessions/:id
func (a *AuditSessionHandler) GetAuditSessionByID(c *gin.Context) {
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := a.service.GetAuditSessionByID(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// AddProductToAudit POST /api/superadmin/audit-sessions/products
func (a *AuditSessionHandler) AddProductToAudit(c *gin.Context) {
	//check request
	var req request.AddProductToAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := a.service.AddProductToAudit(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// RemoveProductFromAudit DELETE /api/superadmin/audit-sessions/:sessionId/products/:productId
func (a *AuditSessionHandler) RemoveProductFromAudit(c *gin.Context) {
	//get sessionID/productID from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	productID, err := getIDParam(c, "productId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = a.service.RemoveProductFromAudit(sessionID, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully removed product from audit"})
}

// CloseAuditSession PATCH /api/superadmin/audit-sessions/:id/close
func (a *AuditSessionHandler) CloseAuditSession(c *gin.Context) {
	//get id from url
	id, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = a.service.CloseAuditSession(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully closed audit session"})
}

// UpdateAuditSession PUT /api/superadmin/audit-sessions/:id
func (a *AuditSessionHandler) UpdateAuditSession(c *gin.Context) {
	//get id from URL
	id, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//bind request
	var req request.UpdateAuditSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	err = a.service.UpdateAuditSession(id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated audit session"})
}
