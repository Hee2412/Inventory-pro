package domain

import (
	"time"
)

type OrderSession struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Title        string    `gorm:"size:255;not null" json:"title"`
	OrderCycle   string    `gorm:"size:255;not null" json:"order_cycle"`
	Status       string    `gorm:"size:50" json:"status"`
	Deadline     time.Time `json:"deadline"`
	DeliveryDate time.Time `gorm:"type:date" json:"delivery_date"`

	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (OrderSession) TableName() string {
	return "order_sessions"
}
