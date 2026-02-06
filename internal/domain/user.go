package domain

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint   `gorm:"primary_key" json:"id"`
	Username  string `gorm:"unique;not null;size:255" json:"username"`
	Password  string `gorm:"not null;size:255" json:"-"`
	Role      string `gorm:"not null;size:255" json:"role"`
	StoreName string `gorm:"not null;size:255" json:"store_name"`
	StoreCode string `gorm:"unique;size:255" json:"store_code"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`

	CreatedBy *uint          `json:"created_by,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	LastLogin *time.Time     `json:"last_login_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Role == "user" && u.StoreCode == "" {
		u.StoreCode = fmt.Sprintf("TEMP-%d", time.Now().UnixNano())
	}
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.Role == "user" {
		newCode := fmt.Sprintf("ST%06d", u.ID)
		return tx.Model(u).Update("store_code", newCode).Error
	}
	return nil
}
