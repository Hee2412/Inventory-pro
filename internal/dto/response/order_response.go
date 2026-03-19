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
	Order *StoreOrderResponse `json:"order"`
	Items []OrderItemResponse `json:"items"`
}

type OrderSessionDetailResponse struct {
	Session  *OrderSessionResponse `json:"order"`
	Products []*ProductResponse    `json:"items"`
}

type AdminOrderInSessionResponse struct {
	ID          uint       `json:"id"`
	StoreID     uint       `json:"store_id"`
	StoreName   string     `json:"store_name"`
	Status      string     `json:"status"`
	SubmittedAt *time.Time `json:"submitted_at"`
	CreatedAt   time.Time  `json:"created_at"`
}

type AddProductResponse struct {
	Added   int      `json:"added"`
	Skipped int      `json:"skipped"`
	Errs    []string `json:"errs,omitempty"`
}

type OrderResponse struct {
	ID        uint   `json:"id"`
	StoreID   uint   `json:"store_id"`
	StoreName string `json:"store_name"`
	Status    string `json:"status"`
	SessionID uint   `json:"session_id"`
}

type StoreWithoutOrderResponse struct {
	SessionID         uint                     `json:"session_id"`
	SessionName       string                   `json:"session_name"`
	NotOrdered        int                      `json:"not_order"`
	StoreWithoutOrder []*StoreTrackingResponse `json:"store_without_order"`
}

type StoreTrackingResponse struct {
	StoreID   uint   `json:"store_id"`
	StoreName string `json:"store_name"`
}

type InventoryUpdateResponse struct {
	OrderID uint                   `json:"order_id"`
	StoreID uint                   `json:"store_id"`
	Updated int                    `json:"updated"`
	Items   []InventoryUpdatedItem `json:"items"`
}

type InventoryUpdatedItem struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	AddedQty    float64 `json:"added_qty"`
	NewTotal    float64 `json:"new_total"`
}

type TransferOrderResponse struct {
	ID            uint       `json:"id"`
	FromStoreID   uint       `json:"from_store_id"`
	FromStoreName string     `json:"from_store_name"`
	ToStoreID     uint       `json:"to_store_id"`
	ToStoreName   string     `json:"to_store_name"`
	Status        string     `json:"status"`
	Note          string     `json:"note"`
	CreatedBy     uint       `json:"created_by"`
	CreatedAt     time.Time  `json:"created_at"`
	ApprovedAt    *time.Time `json:"approved_at,omitempty"`
	CancelledAt   *time.Time `json:"cancelled_at,omitempty"`
	CancelReason  string     `json:"cancel_reason,omitempty"`
}

type TransferOrderDetailResponse struct {
	Order *TransferOrderResponse `json:"order"`
	Items []TransferItemResponse `json:"items"`
}

type TransferItemResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	ProductCode string  `json:"product_code"`
	Quantity    float64 `json:"quantity"`
}
