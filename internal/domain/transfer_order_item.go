package domain

type TransferOrderItem struct {
	ID              uint    `gorm:"primary_key" json:"id"`
	TransferOrderID uint    `gorm:"not null" json:"transfer_order_id"`
	ProductID       uint    `gorm:"not null" json:"product_id"`
	ProductName     string  `json:"product_name"`
	ProductCode     string  `json:"product_code"`
	Quantity        float64 `gorm:"type:decimal(10,2)" json:"quantity"`
}

func (TransferOrderItem) TableName() string {
	return "transfer_order_items"
}
