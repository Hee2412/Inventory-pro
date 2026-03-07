package handler

import (
	"Inventory-pro/internal/dto/request"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get storeID from JWT
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get storeID from JWT
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	storeID := userID.(uint)
	//bind request
	var req request.UpdateAuditItemsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//call service
	err = s.service.UpdateAuditItem(sessionID, storeID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fails to update item",
			"error":   err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully updated audit item"})
}

// GetMyAuditReport GET api/store/audit-reports
func (s *StoreAuditHandler) GetMyAuditReport(c *gin.Context) {
	//get storeId from JWT
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	storeID := userID.(uint)
	//call service
	result, err := s.service.GetMyAuditReports(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Fails to get audit report",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}
