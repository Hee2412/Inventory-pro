package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
)

type TransferOrderItemRepository interface {
	Create(items ...*domain.TransferOrderItem) error
	FindByTransferOrderId(transferOrderID uint) ([]*domain.TransferOrderItem, error)
	DeleteByTransferOrderId(transferOrderID uint) error
}
type transferOrderItemRepository struct {
	db *gorm.DB
}

func NewTransferOrderItemRepository(db *gorm.DB) TransferOrderItemRepository {
	return &transferOrderItemRepository{db: db}
}

func (r *transferOrderItemRepository) Create(items ...*domain.TransferOrderItem) error {
	return r.db.Create(items).Error
}

func (r *transferOrderItemRepository) FindByTransferOrderId(transferOrderID uint) ([]*domain.TransferOrderItem, error) {
	var items []*domain.TransferOrderItem
	err := r.db.Where("transfer_order_id = ?", transferOrderID).Find(&items).Error
	return items, err
}

func (r *transferOrderItemRepository) DeleteByTransferOrderId(transferOrderID uint) error {
	return r.db.Where("transfer_order_id = ?", transferOrderID).
		Delete(&domain.TransferOrderItem{}).Error
}
