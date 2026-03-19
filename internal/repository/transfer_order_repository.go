package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type TransferOrderRepository interface {
	Create(order *domain.TransferOrder) error
	FindById(id uint) (*domain.TransferOrder, error)
	Update(order *domain.TransferOrder) error
	FindByFromStore(storeID uint) ([]*domain.TransferOrder, error)
	FindByToStore(storeID uint) ([]*domain.TransferOrder, error)
	FindAllPaginated(page, limit int, status string) ([]*domain.TransferOrder, int64, error)
}

type transferOrderRepository struct {
	db *gorm.DB
}

func NewTransferOrderRepository(db *gorm.DB) TransferOrderRepository {
	return &transferOrderRepository{db: db}
}

func (r *transferOrderRepository) Create(order *domain.TransferOrder) error {
	return r.db.Create(order).Error
}

func (r *transferOrderRepository) FindById(id uint) (*domain.TransferOrder, error) {
	var order domain.TransferOrder
	err := r.db.Where("id = ?", id).First(&order).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *transferOrderRepository) Update(order *domain.TransferOrder) error {
	return r.db.Save(order).Error
}

func (r *transferOrderRepository) FindByFromStore(storeID uint) ([]*domain.TransferOrder, error) {
	var orders []*domain.TransferOrder
	err := r.db.Where("from_store_id = ?", storeID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *transferOrderRepository) FindByToStore(storeID uint) ([]*domain.TransferOrder, error) {
	var orders []*domain.TransferOrder
	err := r.db.Where("to_store_id = ?", storeID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *transferOrderRepository) FindAllPaginated(page, limit int, status string) ([]*domain.TransferOrder, int64, error) {
	var orders []*domain.TransferOrder
	var total int64
	query := r.db.Model(&domain.TransferOrder{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	query.Count(&total)
	offset := (page - 1) * limit
	err := query.Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, total, err
}
