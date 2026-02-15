package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type StoreOrderRepository interface {
	Create(order *domain.StoreOrder) error
	FindById(id uint) (*domain.StoreOrder, error)
	Update(order *domain.StoreOrder) error

	FindByStoreAndSession(sessionId uint, storeId uint) (*domain.StoreOrder, error)
	FindByStoreID(storeId uint) ([]*domain.StoreOrder, error)
	FindBySessionID(sessionID uint) ([]*domain.StoreOrder, error)
	FindByStatus(status string) ([]*domain.StoreOrder, error)
}
type storeOrderRepository struct {
	db *gorm.DB
}

func NewStoreOrderRepository(db *gorm.DB) StoreOrderRepository {
	return &storeOrderRepository{db: db}
}

func (s *storeOrderRepository) Create(order *domain.StoreOrder) error {
	return s.db.Create(order).Error
}

func (s *storeOrderRepository) FindById(id uint) (*domain.StoreOrder, error) {
	var order domain.StoreOrder
	err := s.db.Where("id = ?", id).First(&order).Error
	return &order, err
}

func (s *storeOrderRepository) Update(order *domain.StoreOrder) error {
	return s.db.Save(order).Error
}

func (s *storeOrderRepository) FindByStoreAndSession(sessionId uint, storeId uint) (*domain.StoreOrder, error) {
	var order domain.StoreOrder
	err := s.db.Where("store_id = ? AND session_id = ?", storeId, sessionId).First(&order).Error
	return &order, err
}

func (s *storeOrderRepository) FindByStoreID(storeId uint) ([]*domain.StoreOrder, error) {
	var orders []*domain.StoreOrder
	err := s.db.Where("store_id = ?", storeId).Find(&orders).Error
	return orders, err
}

func (s *storeOrderRepository) FindBySessionID(sessionID uint) ([]*domain.StoreOrder, error) {
	var orders []*domain.StoreOrder
	err := s.db.Where("session_id = ?", sessionID).Find(&orders).Error
	return orders, err
}

func (s *storeOrderRepository) FindByStatus(status string) ([]*domain.StoreOrder, error) {
	var orders []*domain.StoreOrder
	err := s.db.Where("status = ?", status).Find(&orders).Error
	return orders, err
}
