package domain

import "time"

type StoreAuditReport struct {
	ID          uint      `gorm:"primary_key" json:"id"`
	Session     uint      `json:"session_id"`
	Store       uint      `son:"store_id"`
	Product     uint      `json:"product_id"`
	SystemStock float32   `json:"system_stock"`
	ActualStock float32   `json:"actual_stock"`
	Variance    float32   `json:"variance"`
	Status      string    `gorm:"size:50" json:"status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (StoreAuditReport) TableName() string {
	return "audit_reports"
}
