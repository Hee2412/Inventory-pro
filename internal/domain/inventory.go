package domain

import "time"

type StoreInventory struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	StoreID   uint      `gorm:"not null;index:idx_store_product" json:"store_id"`
	ProductID uint      `gorm:"not null;index:idx_store_product" json:"product_id"`
	Quantity  float64   `gorm:"type:decimal(10,2);default:0" json:"quantity"`
	UpdatedBy uint      `json:"updated_by"`
	UpdatedAt time.Time `json:"updated_at"`

	Store   *User    `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	Product *Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}
