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
	ProductID      []uint `json:"product_ids" binding:"required,min=1"`
}

type SubmitAuditReportRequest struct {
	AuditSessionID uint              `json:"audit_session_id" binding:"required"`
	StoreID        uint              `json:"store_id" binding:"required"`
	Items          []AuditReportItem `json:"items" binding:"required,min=1"`
}

type AuditReportItem struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	ActualStock float64 `json:"actual_stock"`
}

type UpdateAuditSessionRequest struct {
	Title     *string    `json:"title"`
	AuditType *string    `json:"audit_type"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	Status    *string    `json:"status" `
}

type UpdateAuditItem struct {
	ProductID   uint    `json:"product_id" binding:"required"`
	ActualStock float64 `json:"actual_stock"`
}

type UpdateAuditItemsRequest struct {
	Items []UpdateAuditItem `json:"items"`
}
type DeclineReportRequest struct {
	Reason string `json:"reason"`
}

type AuditSessionSearchParams struct {
	Status   string     `form:"status"`
	FromDate *time.Time `form:"from_date" time_format:"2006-01-02"`
	ToDate   *time.Time `form:"to_date" time_format:"2006-01-02"`
	Page     int        `form:"page" binding:"omitempty,min=1"`
	Limit    int        `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (p *AuditSessionSearchParams) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}
}
