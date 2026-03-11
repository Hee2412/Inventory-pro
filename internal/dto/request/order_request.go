package request

import "time"

type CreateOrderSessionRequest struct {
	Title        string    `json:"title" binding:"required"`
	OrderCycle   string    `json:"order_cycle" binding:"required"`
	Deadline     time.Time `json:"deadline" binding:"required"`
	DeliveryDate time.Time `json:"delivery_date" binding:"required"`
}

type AddProductToSessionRequest struct {
	SessionID uint   `json:"session_id" binding:"required"`
	ProductID []uint `json:"product_id" binding:"required"`
}

type UpdateOrderItemRequest struct {
	Items []UpdateOrderItem `json:"items" binding:"required,min=1"`
}
type UpdateOrderItem struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"min=0"`
}

type DeclineOrderRequest struct {
	Reason string `json:"reason" binding:"required,min=5"`
}

type OrderSearchParams struct {
	StoreID   *uint      `form:"store_id"`
	SessionID *uint      `form:"session_id"`
	Status    string     `form:"status"`
	FromDate  *time.Time `form:"from_date" time_format:"2006-01-02"`
	ToDate    *time.Time `form:"to_date" time_format:"2006-01-02"`
	Page      int        `form:"page" binding:"omitempty,min=1"`
	Limit     int        `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (p *OrderSearchParams) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}
}

type SessionSearchParams struct {
	Status   string     `form:"status"`
	FromDate *time.Time `form:"from_date" time_format:"2006-01-02"`
	ToDate   *time.Time `form:"to_date" time_format:"2006-01-02"`
	Page     int        `form:"page" binding:"omitempty,min=1"`
	Limit    int        `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (p *SessionSearchParams) SetDefaults() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 20
	}
}
