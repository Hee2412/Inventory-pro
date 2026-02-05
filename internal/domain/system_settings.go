package domain

import "time"

type SystemSettings struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Key       string    `gorm:"unique;size:100" json:"key"` //order_auto_lock, max_draft_hours
	Value     string    `json:"value"`
	Desc      string    `json:"desc"`
	UpdatedBy *uint     `json:"updated_by,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (SystemSettings) TableName() string {
	return "system_settings"
}
