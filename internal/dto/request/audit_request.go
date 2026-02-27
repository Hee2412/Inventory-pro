package request

import (
	"time"
)

type CreateAuditSessionRequest struct {
	Title     string    `json:"title"`
	AuditType string    `json:"audit_type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

type AddProductToAuditRequest struct {
	AuditSessionID uint   `json:"audit_session_id" binding:"required"`
	ProductID      []uint `json:"product_id" binding:"required"`
}

type SubmitAuditReportRequest struct {
	AuditSessionID uint `json:"audit_session_id" binding:"required"`
	StoreID        uint `json:"store_id" binding:"required"`
	Items          []struct {
		ProductID   uint    `json:"product_id" binding:"required"`
		ActualStock float32 `json:"actual_stock"`
	} `json:"items" binding:"required"`
}

type UpdateAuditItemRequest struct {
	AuditSessionID uint   `json:"audit_session_id" `
	Title          string `json:"title"`
	AuditType      string `json:"audit_type"`
	EndDate        string `json:"end_date"`
	Status         string `json:"status" `
}
