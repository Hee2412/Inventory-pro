package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type StoreAuditHandler struct {
	service service.StoreAuditService
}

func NewStoreAuditHandler(storeAuditHandler service.StoreAuditService) *StoreAuditHandler {
	return &StoreAuditHandler{service: storeAuditHandler}
}

// GetAuditReport GET /api/store/audit-sessions/:sessionId/report
func (s *StoreAuditHandler) GetAuditReport(c *gin.Context) {
	//get sessionID from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//get storeID from JWT
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "Invalid token")
		return
	}
	storeID := userID.(uint)
	//call service
	report, err := s.service.GetAuditReport(sessionID, storeID)
	c.JSON(http.StatusOK, gin.H{"data": report})
}

// UpdateAuditItem PUT /api/store/audit-sessions/:sessionId/items
func (s *StoreAuditHandler) UpdateAuditItem(c *gin.Context) {
	//get sessionID from URL
	sessionID, err := getIDParam(c, "sessionId")
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//get storeID from JWT
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "Invalid token")
		return
	}
	storeID := userID.(uint)
	//bind request
	var req request.UpdateAuditItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	//call service
	err = s.service.UpdateAuditItem(sessionID, storeID, req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Message(c, "Updated AuditItem")
}

// GetMyAuditReport GET api/store/audit-reports
func (s *StoreAuditHandler) GetMyAuditReport(c *gin.Context) {
	//get storeId from JWT
	userID, exists := c.Get("userId")
	if !exists {
		response.Unauthorized(c, "Invalid token")
		return
	}
	storeID := userID.(uint)
	//call service
	result, err := s.service.GetMyAuditReports(storeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, result)
}
