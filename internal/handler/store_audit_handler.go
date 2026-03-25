package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
)

type StoreAuditHandler struct {
	service service.StoreAuditService
}

func NewStoreAuditHandler(storeAuditHandler service.StoreAuditService) *StoreAuditHandler {
	return &StoreAuditHandler{service: storeAuditHandler}
}

// GetAuditReport GET /api/store/audit-sessions/:sessionId/report
func (s *StoreAuditHandler) GetAuditReport(c *gin.Context) {
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	storeID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	report, err := s.service.GetAuditReport(sessionID, storeID)
	response.Success(c, report)
}

// UpdateAuditItem PUT /api/store/audit-sessions/:sessionId/items
func (s *StoreAuditHandler) UpdateAuditItem(c *gin.Context) {
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	storeID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	var req request.UpdateAuditItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	err = s.service.UpdateAuditItem(sessionID, storeID, req)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Updated AuditItem")
}

// GetMyAuditReport GET api/store/audit-reports
func (s *StoreAuditHandler) GetMyAuditReport(c *gin.Context) {
	storeID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := s.service.GetMyAuditReports(storeID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}
