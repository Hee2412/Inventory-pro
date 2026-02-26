package request

type UpdateProductRequest struct {
	ProductName *string  `json:"product_name"`
	Unit        *string  `json:"unit"`
	MOQ         *float32 `json:"moq"`
	OM          *float32 `json:"om"`
	Type        *string  `json:"type"`
	OrderCycle  *string  `json:"order_cycle"`
	AuditCycle  *string  `json:"audit_cycle"`
	IsActive    *bool    `json:"is_active"`
}

type CreateProductRequest struct {
	ProductName string  `json:"product_name" binding:"required"`
	Unit        string  `json:"unit" binding:"required"`
	MOQ         float32 `json:"moq" binding:"required"`
	OM          float32 `json:"om" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	OrderCycle  string  `json:"order_cycle" binding:"required"`
	AuditCycle  string  `json:"audit_cycle" binding:"required"`
}
