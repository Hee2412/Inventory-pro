package response

import "time"

type InventoryResponse struct {
	ID          uint      `json:"id"`
	StoreID     uint      `json:"store_id"`
	StoreName   string    `json:"store_name"`
	ProductID   uint      `json:"product_id"`
	ProductName string    `json:"product_name"`
	ProductCode string    `json:"product_code"`
	Quantity    float64   `json:"quantity"`
	UpdatedBy   uint      `json:"updated_by"`
	UpdatedAt   time.Time `json:"updated_at"`
}
