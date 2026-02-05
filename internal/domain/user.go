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
	StoreCode string `gorm:"unique;not null;size:255" json:"store_code"`
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
	if u.Role == "store" && u.StoreCode == "" {

	}
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) (err error) {
	if u.Role == "store" && u.StoreCode == "" {
		storeCode := fmt.Sprintf("ST%06d", u.ID)
		tx.Model(u).Update("store_code", storeCode)
	}
	return nil
}
