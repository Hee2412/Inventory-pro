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
	SessionInfo AuditSessionResponse       `json:"session_info"`
	Report      []StoreAuditReportResponse `json:"report"`
}

type StoreAuditReportResponse struct {
	StoreId     uint    `json:"store_id"`
	StoreName   string  `json:"store_name"`
	ProductName string  `json:"product_name"`
	SystemStock float32 `json:"system_stock"`
	ActualStock float32 `json:"actual_stock"`
	Variance    float32 `json:"variance"`
	Status      string  `json:"status"`
}
type AuditReportDetailsResponse struct {
	SessionTitle string                     `json:"session_title"`
	StoreName    string                     `json:"store_name"`
	TotalItems   int                        `json:"total_items"`
	TotalOverlap int                        `json:"total_overlap"`
	SubmittedAt  time.Time                  `json:"submitted_at"`
	Items        []StoreAuditReportResponse `json:"items"`
}
