package domain

type OrderSessionProduct struct {
	ID        uint `gorm:"primary_key" json:"id"`
	SessionID uint `json:"session_id"`
	ProductID uint `json:"product_id"`
}

func (OrderSessionProduct) TableName() string {
	return "order_session_products"
}
