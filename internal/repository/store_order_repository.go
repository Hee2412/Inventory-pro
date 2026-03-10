package repository

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/pkg/pagination"
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
	FindAllPaginated(page, limit int) ([]*domain.StoreOrder, int64, error)
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
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *storeOrderRepository) Update(order *domain.StoreOrder) error {
	return s.db.Save(order).Error
}

func (s *storeOrderRepository) FindByStoreAndSession(sessionId uint, storeId uint) (*domain.StoreOrder, error) {
	var order domain.StoreOrder
	err := s.db.Where("store_id = ? AND session_id = ?", storeId, sessionId).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *storeOrderRepository) FindByStoreID(storeId uint) ([]*domain.StoreOrder, error) {
	var orders []*domain.StoreOrder
	err := s.db.Where("store_id = ?", storeId).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *storeOrderRepository) FindBySessionID(sessionID uint) ([]*domain.StoreOrder, error) {
	var orders []*domain.StoreOrder
	err := s.db.Where("session_id = ?", sessionID).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *storeOrderRepository) FindByStatus(status string) ([]*domain.StoreOrder, error) {
	var orders []*domain.StoreOrder
	err := s.db.Where("status = ?", status).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *storeOrderRepository) FindAllPaginated(page, limit int) ([]*domain.StoreOrder, int64, error) {
	var users []*domain.StoreOrder
	var total int64
	// count total
	if err := s.db.Model(&domain.StoreOrder{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	//get paginated data
	err := s.db.Scopes(pagination.Paginate(page, limit)).
		Find(&users).Error

	return users, total, err
}
