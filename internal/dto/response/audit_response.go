package response

import "time"

type AuditSessionResponse struct {
	SessionID uint      `json:"session_id"`
	Title     string    `json:"title"`
	AuditType string    `json:"audit_type"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	Status    string    `json:"status"`
}

type AuditSessionDetailsResponse struct {
	SessionInfo AuditSessionResponse        `json:"session_info"`
	Report      []StoreAuditSummaryResponse `json:"report"`
}

type StoreAuditSummaryResponse struct {
	StoreId   uint   `json:"store_id"`
	StoreName string `json:"store_name"`
	Status    string `json:"status"`
}
type AuditReportItemResponse struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	SystemStock float64 `json:"system_stock"`
	ActualStock float64 `json:"actual_stock"`
	Variance    float64 `json:"variance"`
}
type AuditReportDetailResponse struct {
	SessionTitle  string                    `json:"session_title"`
	StoreName     string                    `json:"store_name"`
	StoreID       uint                      `json:"store_id"`
	TotalItems    int                       `json:"total_items"`
	TotalVariance float64                   `json:"total_variance"`
	SubmittedAt   *time.Time                `json:"submitted_at,omitempty"`
	Status        string                    `json:"status"`
	Items         []AuditReportItemResponse `json:"items"`
}
