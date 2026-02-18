package request

import "time"

type CreateOrderSessionRequest struct {
	Title        string    `json:"title" binding:"required"`
	OrderCycle   string    `json:"order_cycle" binding:"required"`
	Deadline     time.Time `json:"deadline" binding:"required"`
	DeliveryDate time.Time `json:"delivery_date" binding:"required"`
}

type AddProductToSessionRequest struct {
	SessionID uint `json:"session_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
}

type UpdateOrderItemRequest struct {
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  float64 `json:"quantity" binding:"required"`
}
