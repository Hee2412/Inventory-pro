package domain

import "time"

type OrderItems struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	OrderID     uint      `json:"order_id"`
	ProductID   uint      `json:"product_id"`
	ProductName string    `json:"product_name"`
	ProductCode string    `json:"product_code"`
	Quantity    float64   `gorm:"default:0" json:"quantity"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (OrderItems) TableName() string {
	return "store_order_items"
}
