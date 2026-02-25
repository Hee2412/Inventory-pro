package response

import "time"

type OrderSessionResponse struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	OrderCycle string    `json:"order_cycle"`
	Status     string    `json:"status"`
	Deadline   time.Time `json:"deadline"`
	DeliveryAt time.Time `json:"delivery_at"`
	CreatedAt  time.Time `json:"created_at"`
}
type StoreOrderResponse struct {
	ID          uint       `json:"id"`
	SessionID   uint       `json:"session_id"`
	StoreID     uint       `json:"store_id"`
	Status      string     `json:"status"`
	SubmittedAt *time.Time `json:"submitted_at"`
	ApprovedAt  *time.Time `json:"approved_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type OrderItemResponse struct {
	ID          uint    `json:"id"`
	OrderID     uint    `json:"order_id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	ProductCode string  `json:"product_code"`
	Quantity    float64 `json:"quantity"`
}

type StoreOrderDetailResponse struct {
	Order StoreOrderResponse  `json:"order"`
	Items []OrderItemResponse `json:"items"`
}

type OrderSessionDetailResponse struct {
	Session  OrderSessionResponse `json:"order"`
	Products []ProductResponse    `json:"items"`
}

type AdminOrderInSessionResponse struct {
	ID          uint       `json:"id"`
	StoreID     uint       `json:"store_id"`
	StoreName   string     `json:"store_name"`
	Status      string     `json:"status"`
	SubmittedAt *time.Time `json:"submitted_at"`
	CreatedAt   time.Time  `json:"created_at"`
}
