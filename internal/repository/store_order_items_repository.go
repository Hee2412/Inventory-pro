package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type StoreOrderItemRepository interface {
	Create(order *domain.OrderItems) error
	FindById(id uint) (*domain.OrderItems, error)
	Update(order *domain.OrderItems) error
	Delete(id uint) error
	FindByOrderId(orderId uint) ([]domain.OrderItems, error)
	DeleteByOrderId(orderId uint) error
}

type storeOrderItems struct {
	db *gorm.DB
}

func NewStoreOrderItems(db *gorm.DB) StoreOrderItemRepository {
	return &storeOrderItems{db: db}
}

func (s *storeOrderItems) Create(order *domain.OrderItems) error {
	return s.db.Create(order).Error
}

func (s *storeOrderItems) FindById(id uint) (*domain.OrderItems, error) {
	var order domain.OrderItems
	err := s.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *storeOrderItems) Update(order *domain.OrderItems) error {
	return s.db.Save(order).Error
}

func (s *storeOrderItems) Delete(id uint) error {
	return s.db.Delete(&domain.OrderItems{}, id).Error
}

func (s *storeOrderItems) FindByOrderId(orderId uint) ([]domain.OrderItems, error) {
	var order []domain.OrderItems
	err := s.db.Where("order_id = ?", orderId).Find(&order).Error
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *storeOrderItems) DeleteByOrderId(orderId uint) error {
	return s.db.Where("order_id = ?", orderId).Delete(&domain.OrderItems{}).Error
}
