package domain

import (
	"time"
)

type StoreAuditReport struct {
	ID           uint         `gorm:"primary_key" json:"id"`
	SessionID    uint         `json:"session_id"`
	StoreID      uint         `json:"store_id"`
	ProductID    uint         `json:"product_id"`
	ProductName  string       `json:"product_name"`
	StoreName    string       `json:"store_name"`
	Product      Product      `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Store        *User        `gorm:"foreignKey:StoreID" json:"store,omitempty"`
	OrderSession OrderSession `gorm:"foreignKey:SessionID" json:"order_session,omitempty"`
	SystemStock  float64      `json:"system_stock"`
	ActualStock  float64      `json:"actual_stock"`
	Variance     float64      `json:"variance"`
	Status       string       `gorm:"size:50" json:"status"`
	UpdatedAt    time.Time    `json:"updated_at"`
	ApprovedAt   *time.Time   `json:"approved_at"`
	ApprovedBy   uint         `json:"approved_by"`
}

func (StoreAuditReport) TableName() string {
	return "audit_reports"
}
