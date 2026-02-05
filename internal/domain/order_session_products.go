package domain

type OrderSessionProducts struct {
	ID        uint `gorm:"primary_key" json:"id"`
	SessionID uint `json:"session_id"`
	ProductID uint `json:"product_id"`
}

func (OrderSessionProducts) TableName() string {
	return "order_session_products"
}
