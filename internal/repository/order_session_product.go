package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type OrderSessionProducts interface {
	Create(order *domain.OrderSessionProducts) error
	FindBySessionId(sessionId uint) (*domain.OrderSessionProducts, error)
	Delete(id uint) error
	FindBySessionAndProduct(sessionID uint, productID uint) (*domain.OrderSessionProducts, error)
}

type orderSessionProducts struct {
	db *gorm.DB
}

func NewOrderSessionProducts(db *gorm.DB) OrderSessionProducts {
	return &orderSessionProducts{db: db}
}

func (o *orderSessionProducts) Create(order *domain.OrderSessionProducts) error {
	return o.db.Create(order).Error
}

func (o *orderSessionProducts) FindBySessionId(sessionId uint) (*domain.OrderSessionProducts, error) {
	var order domain.OrderSessionProducts
	err := o.db.Where("session_id = ?", sessionId).First(&order).Error
	return &order, err
}

func (o *orderSessionProducts) Delete(id uint) error {
	return o.db.Delete(&domain.OrderSession{}, id).Error
}

func (o *orderSessionProducts) FindBySessionAndProduct(sessionID uint, productID uint) (*domain.OrderSessionProducts, error) {
	var order domain.OrderSessionProducts
	err := o.db.Where("session_id = ? and product_id = ?", sessionID, productID).First(&order).Error
	return &order, err
}
