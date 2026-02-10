package domain

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	ProductName string  `gorm:"size:255;not null" json:"product_name"`
	ProductCode string  `gorm:"size:100;unique" json:"product_code"`
	Unit        string  `gorm:"size:50" json:"unit"`
	MOQ         float32 `json:"moq"`
	OM          float32 `json:"om"`
	Type        string  `gorm:"size:50;not null" json:"type"`
	OrderCycle  string  `gorm:"size:50;not null" json:"order_cycle"`
	AuditCycle  string  `gorm:"size:50;not null" json:"audit_cycle"`
	IsActive    bool    `gorm:"default:true" json:"is_active"`

	CreatedBy *uint      `json:"created_by,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	UpdateBy  *uint      `json:"update_by,omitempty"`
	DeletedAt *time.Time `gorm:"index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}

func (p *Product) AfterCreate(tx *gorm.DB) (err error) {
	if p.ProductCode == "" {
		productCode := fmt.Sprintf("PRD%06d", p.ID)
		tx.Model(p).Update("product_code", productCode)
	}
	return nil
}
