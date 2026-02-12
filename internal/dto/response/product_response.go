package response

type ProductResponse struct {
	ID          uint    `json:"id"`
	ProductName string  `json:"product_name"`
	ProductCode string  `json:"product_code"`
	Unit        string  `json:"unit"`
	MOQ         float32 `json:"moq"`
	OM          float32 `json:"om"`
	Type        string  `json:"type"`
	OrderCycle  string  `json:"order_cycle"`
	AuditCycle  string  `json:"audit_cycle"`
	IsActive    bool    `json:"is_active"`
}
