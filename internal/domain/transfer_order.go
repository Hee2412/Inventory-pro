package domain

import "time"

type TransferOrder struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	FromStoreID  uint       `gorm:"not null" json:"from_store_id"`
	ToStoreID    uint       `gorm:"not null" json:"to_store_id"`
	Status       string     `gorm:"size:50;default:'PENDING'" json:"status"`
	Note         string     `json:"note"`
	CreatedBy    uint       `gorm:"not null" json:"created_by"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ApprovedAt   *time.Time `json:"approved_at,omitempty"`
	ApprovedBy   *uint      `json:"approved_by,omitempty"`
	CancelledAt  *time.Time `json:"cancelled_at,omitempty"`
	CancelReason string     `json:"cancel_reason,omitempty"`
	CancelBy     *uint      `json:"cancel_by,omitempty"`
}

func (TransferOrder) TableName() string {
	return "transfer_orders"
}
