package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type StoreOrderItems interface {
	Create(order *domain.OrderItems) error
	FindById(id int) (*domain.OrderItems, error)
	Update(order *domain.OrderItems) error
	Delete(id int) error
}

type storeOrderItems struct {
	db *gorm.DB
}

func NewStoreOrderItems(db *gorm.DB) StoreOrderItems {
	return &storeOrderItems{db: db}
}

func (s *storeOrderItems) Create(order *domain.OrderItems) error {
	return s.db.Create(order).Error
}

func (s *storeOrderItems) FindById(id int) (*domain.OrderItems, error) {
	var order domain.OrderItems
	return &order, s.db.Where("id = ?", id).First(&order).Error
}

func (s *storeOrderItems) Update(order *domain.OrderItems) error {
	return s.db.Save(order).Error
}

func (s *storeOrderItems) Delete(id int) error {
	return s.db.Where("id = ?", id).Delete(&domain.OrderItems{}).Error
}
