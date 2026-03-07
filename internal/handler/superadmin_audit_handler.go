package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SuperadminAuditHandler struct {
	service service.SuperAdminAuditService
}

func NewSuperadminAuditHandler(superadminAuditHandler service.SuperAdminAuditService) *SuperadminAuditHandler {
	return &SuperadminAuditHandler{service: superadminAuditHandler}
}

// GetAllReportsInSession GET /api/superadmin/audit-sessions/:sessionId/reports
func (s *SuperadminAuditHandler) GetAllReportsInSession(c *gin.Context) {
	//get sessionID from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// call service
	result, err := s.service.GetAllReportsInSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetReportDetail GET /api/superadmin/audit-sessions/:sessionId/stores/:storeId
func (s *SuperadminAuditHandler) GetReportDetail(c *gin.Context) {
	//get sessionID/storeId from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	storeID, err := getIDParam(c, "storeId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	result, err := s.service.GetReportDetail(sessionID, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// GetAuditSummary GET /api/superadmin/audit-sessions/:sessionId/summary
func (s *SuperadminAuditHandler) GetAuditSummary(c *gin.Context) {
	//get sessionID from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	result, err := s.service.GetAuditSummary(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

// ApproveStoreReport POST /api/superadmin/audit-sessions/:sessionId/stores/:storeId/approve
func (s *SuperadminAuditHandler) ApproveStoreReport(c *gin.Context) {
	//get sessionID/storeID from URL
	storeID, err := getIDParam(c, "storeId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get adminID from JWT
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	adminID := userID.(uint)
	//call service
	err = s.service.ApproveStoreReport(storeID, sessionID, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Store report approved successfully"})
}

// DeclineStoreReport POST /api/superadmin/audit-sessions/:sessionId/stores/:storeId/decline
func (s *SuperadminAuditHandler) DeclineStoreReport(c *gin.Context) {
	//get sessionID/storeID from URL
	storeID, err := getIDParam(c, "storeId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get adminID from JWT
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	adminID := userID.(uint)
	//var reason
	var req request.DeclineReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	err = s.service.DeclineStoreReport(storeID, sessionID, req.Reason, adminID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Store report declined successfully"})
}
