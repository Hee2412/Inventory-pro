package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
)

type SuperadminAuditHandler struct {
	service service.SuperAdminAuditService
}

func NewSuperadminAuditHandler(superadminAuditHandler service.SuperAdminAuditService) *SuperadminAuditHandler {
	return &SuperadminAuditHandler{service: superadminAuditHandler}
}

// GetAllReportsInSession GET /api/superadmin/audit-sessions/:sessionId/reports
func (s *SuperadminAuditHandler) GetAllReportsInSession(c *gin.Context) {
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := s.service.GetAllReportsInSession(sessionID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetReportDetail GET /api/superadmin/audit-sessions/:sessionId/stores/:storeId
func (s *SuperadminAuditHandler) GetReportDetail(c *gin.Context) {
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	storeID, err := getIDParam(c, "storeId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := s.service.GetReportDetail(sessionID, storeID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetAuditSummary GET /api/superadmin/audit-sessions/:sessionId/summary
func (s *SuperadminAuditHandler) GetAuditSummary(c *gin.Context) {
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := s.service.GetAuditSummary(sessionID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// ApproveStoreReport POST /api/superadmin/audit-sessions/:sessionId/stores/:storeId/approve
func (s *SuperadminAuditHandler) ApproveStoreReport(c *gin.Context) {
	storeID, err := getIDParam(c, "storeId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	adminID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	err = s.service.ApproveStoreReport(storeID, sessionID, adminID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Report approved")
}

// DeclineStoreReport POST /api/superadmin/audit-sessions/:sessionId/stores/:storeId/decline
func (s *SuperadminAuditHandler) DeclineStoreReport(c *gin.Context) {
	storeID, err := getIDParam(c, "storeId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	adminID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	var req request.DeclineReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	err = s.service.DeclineStoreReport(storeID, sessionID, req.Reason, adminID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Report declined")
}
