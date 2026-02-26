package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type OrderSessionProductRepository interface {
	Create(order *domain.OrderSessionProducts) error
	FindBySessionId(sessionId uint) ([]*domain.OrderSessionProducts, error)
	Delete(id uint) error
	FindBySessionAndProduct(sessionID uint, productID uint) (*domain.OrderSessionProducts, error)
}

type orderSessionProductRepository struct {
	db *gorm.DB
}

func NewOrderSessionProductRepository(db *gorm.DB) OrderSessionProductRepository {
	return &orderSessionProductRepository{db: db}
}

func (o *orderSessionProductRepository) Create(order *domain.OrderSessionProducts) error {
	return o.db.Create(order).Error
}

func (o *orderSessionProductRepository) FindBySessionId(sessionId uint) ([]*domain.OrderSessionProducts, error) {
	var order []*domain.OrderSessionProducts
	err := o.db.Where("session_id = ?", sessionId).Find(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *orderSessionProductRepository) Delete(id uint) error {
	return o.db.Delete(&domain.OrderSessionProducts{}, id).Error
}

func (o *orderSessionProductRepository) FindBySessionAndProduct(sessionID uint, productID uint) (*domain.OrderSessionProducts, error) {
	var order domain.OrderSessionProducts
	err := o.db.Where("session_id = ? and product_id = ?", sessionID, productID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
