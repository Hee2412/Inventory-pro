package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type OrderSessionProductRepository interface {
	Create(order ...*domain.OrderSessionProduct) error
	FindBySessionId(sessionId uint) ([]*domain.OrderSessionProduct, error)
	Delete(id uint) error
	FindBySessionAndProduct(sessionID uint, productID uint) (*domain.OrderSessionProduct, error)
}

type orderSessionProductRepository struct {
	db *gorm.DB
}

func NewOrderSessionProductRepository(db *gorm.DB) OrderSessionProductRepository {
	return &orderSessionProductRepository{db: db}
}

func (o *orderSessionProductRepository) Create(products ...*domain.OrderSessionProduct) error {
	return o.db.Create(products).Error
}

func (o *orderSessionProductRepository) FindBySessionId(sessionId uint) ([]*domain.OrderSessionProduct, error) {
	var order []*domain.OrderSessionProduct
	err := o.db.Where("session_id = ?", sessionId).Find(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *orderSessionProductRepository) Delete(id uint) error {
	return o.db.Delete(&domain.OrderSessionProduct{}, id).Error
}

func (o *orderSessionProductRepository) FindBySessionAndProduct(sessionID uint, productID uint) (*domain.OrderSessionProduct, error) {
	var order domain.OrderSessionProduct
	err := o.db.Where("session_id = ? and product_id = ?", sessionID, productID).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}
