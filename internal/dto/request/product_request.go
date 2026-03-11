package request

type UpdateProductRequest struct {
	ProductName *string  `json:"product_name"`
	Unit        *string  `json:"unit"`
	MOQ         *float64 `json:"moq"`
	OM          *float64 `json:"om"`
	Type        *string  `json:"type"`
	OrderCycle  *string  `json:"order_cycle"`
	AuditCycle  *string  `json:"audit_cycle"`
	IsActive    *bool    `json:"is_active"`
}

type CreateProductRequest struct {
	ProductName string  `json:"product_name" binding:"required"`
	Unit        string  `json:"unit" binding:"required"`
	MOQ         float64 `json:"moq" binding:"required"`
	OM          float64 `json:"om" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	OrderCycle  string  `json:"order_cycle" binding:"required"`
	AuditCycle  string  `json:"audit_cycle" binding:"required"`
}

type ProductSearchParams struct {
	Search   string   `form:"search"`
	IsActive *bool    `form:"is_active"`
	MinPrice *float64 `form:"min_price"`
	MaxPrice *float64 `form:"max_price"`
	Page     int      `form:"page" binding:"omitempty,min=1"`
	Limit    int      `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (p *ProductSearchParams) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}
}
