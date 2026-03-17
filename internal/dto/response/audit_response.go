package response

import (
	"time"
)

type AuditSessionResponse struct {
	SessionID uint      `json:"session_id"`
	Title     string    `json:"title"`
	AuditType string    `json:"audit_type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
	CreatedBy uint      `json:"created_by"`
}

type AuditSessionDetailsResponse struct {
	SessionInfo *AuditSessionResponse       `json:"session_info"`
	Report      []StoreAuditSummaryResponse `json:"report"`
}

type StoreAuditSummaryResponse struct {
	StoreId   uint   `json:"store_id"`
	StoreName string `json:"store_name"`
	Status    string `json:"status"`
}
type AuditItemsResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	SystemStock float64 `json:"system_stock"`
	ActualStock float64 `json:"actual_stock"`
	Variance    float64 `json:"variance"`
}

type AuditReportItemDetailResponse struct {
	SessionTitle string                `json:"session_title"`
	TotalItems   int                   `json:"total_items"`
	Items        []*AuditItemsResponse `json:"items"`
}

type AuditSummaryResponse struct {
	SessionTitle     string       `json:"session_title"`
	TotalStores      int          `json:"total_stores"`
	StoresApproved   int          `json:"stores_approved"`
	StoreDraft       int          `json:"store_draft"`
	TotalProducts    int          `json:"total_products"`
	TotalVariance    float64      `json:"total_variance"`
	StoresWithIssues []StoreIssue `json:"stores_with_issues"`
}

type StoreIssue struct {
	StoreID   uint    `json:"store_id"`
	StoreName string  `json:"store_name"`
	Variance  float64 `json:"variance"`
	Status    string  `json:"status"`
}

type AuditTrackingResponse struct {
	Completed       int              `json:"completed"`
	Incomplete      int              `json:"incomplete"`
	IncompleteStore []*StoreTracking `json:"incomplete_store"`
}

type StoreTracking struct {
	StoreID   uint   `json:"store_id"`
	StoreName string `json:"store_name"`
}
