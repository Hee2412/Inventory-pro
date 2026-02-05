package domain

import "time"

type StoreOrders struct {
	ID          uint       `gorm:"primary_key" json:"id"`
	SessionID   uint       `json:"session_id"`
	StoreID     uint       `json:"store_id"`
	Status      string     `gorm:"size:50" json:"status"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
	ApproveAt   *time.Time `json:"approve_at,omitempty"`
	ApproveBy   *uint      `json:"approve_by,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (StoreOrders) TableName() string {
	return "store_orders"
}
