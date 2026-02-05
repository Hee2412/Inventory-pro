package domain

import (
	"gorm.io/datatypes"
	"time"
)

type SystemLogs struct {
	ID         uint           `gorm:"primary_key;" json:"id"`
	UserID     uint           `json:"user_id"`
	ActionType string         `gorm:"size:100" json:"action_type"`
	TargetType string         `gorm:"size:100" json:"target_type"`
	TargetID   uint           `json:"target_id"`
	OldData    datatypes.JSON `gorm:"type:json" json:"old_data"`
	NewData    datatypes.JSON `gorm:"type:json" json:"new_data"`
	IpAddr     string         `json:"ip_addr"`
	UserAgent  string         `json:"user_agent"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

func (SystemLogs) TableName() string {
	return "system_log"
}
