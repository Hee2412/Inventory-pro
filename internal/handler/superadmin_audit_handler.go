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
	//get sessionID from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	// call service
	result, err := s.service.GetAllReportsInSession(sessionID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetReportDetail GET /api/superadmin/audit-sessions/:sessionId/stores/:storeId
func (s *SuperadminAuditHandler) GetReportDetail(c *gin.Context) {
	//get sessionID/storeId from URL
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
	//call service
	result, err := s.service.GetReportDetail(sessionID, storeID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetAuditSummary GET /api/superadmin/audit-sessions/:sessionId/summary
func (s *SuperadminAuditHandler) GetAuditSummary(c *gin.Context) {
	//get sessionID from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	//call service
	result, err := s.service.GetAuditSummary(sessionID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// ApproveStoreReport POST /api/superadmin/audit-sessions/:sessionId/stores/:storeId/approve
func (s *SuperadminAuditHandler) ApproveStoreReport(c *gin.Context) {
	//get sessionID/storeID from URL
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
	//get adminID from JWT
	adminID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	//call service
	err = s.service.ApproveStoreReport(storeID, sessionID, adminID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Report approved")
}

// DeclineStoreReport POST /api/superadmin/audit-sessions/:sessionId/stores/:storeId/decline
func (s *SuperadminAuditHandler) DeclineStoreReport(c *gin.Context) {
	//get sessionID/storeID from URL
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
	//get adminID from JWT
	adminID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	//var reason
	var req request.DeclineReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	//call service
	err = s.service.DeclineStoreReport(storeID, sessionID, req.Reason, adminID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Report declined")
}
