package repository

import (
	"Inventory-pro/internal/domain"
	"gorm.io/gorm"
	"time"
)

type StoreInventoryRepository interface {
	FindByStoreAndProduct(storeID, productID uint) (*domain.StoreInventory, error)
	FindByStoreID(storeID uint) ([]*domain.StoreInventory, error)
	FindByProductID(productID uint) ([]*domain.StoreInventory, error)
	Upsert(inventory *domain.StoreInventory) error
	UpdateQuantity(storeID, productID uint, quantity float64, updatedBy uint) error
	AdjustQuantity(storeID, productID uint, delta float64, updatedBy uint) error
	GetAllInventory() ([]*domain.StoreInventory, error)
}

type storeInventoryRepository struct {
	db *gorm.DB
}

func NewStoreInventoryRepository(db *gorm.DB) StoreInventoryRepository {
	return &storeInventoryRepository{db: db}
}

func (r *storeInventoryRepository) FindByStoreAndProduct(storeID, productID uint) (*domain.StoreInventory, error) {
	var inventory domain.StoreInventory
	err := r.db.Where("store_id = ? AND product_id = ?", storeID, productID).
		Preload("Store").
		Preload("Product").
		First(&inventory).Error
	return &inventory, err
}

func (r *storeInventoryRepository) FindByStoreID(storeID uint) ([]*domain.StoreInventory, error) {
	var inventories []*domain.StoreInventory
	err := r.db.Where("store_id = ?", storeID).
		Preload("Product").
		Find(&inventories).Error
	return inventories, err
}

func (r *storeInventoryRepository) FindByProductID(productID uint) ([]*domain.StoreInventory, error) {
	var inventories []*domain.StoreInventory
	err := r.db.Where("product_id = ?", productID).
		Preload("Store").
		Find(&inventories).Error
	return inventories, err
}

func (r *storeInventoryRepository) Upsert(inventory *domain.StoreInventory) error {
	return r.db.Save(inventory).Error
}

func (r *storeInventoryRepository) UpdateQuantity(storeID, productID uint, quantity float64, updatedBy uint) error {
	return r.db.Model(&domain.StoreInventory{}).
		Where("store_id = ? AND product_id = ?", storeID, productID).
		Updates(map[string]interface{}{
			"quantity":   quantity,
			"updated_by": updatedBy,
			"updated_at": time.Now(),
		}).Error
}

func (r *storeInventoryRepository) AdjustQuantity(storeID, productID uint, delta float64, updatedBy uint) error {
	return r.db.Exec(`
        INSERT INTO store_inventories (store_id, product_id, quantity, updated_by, updated_at)
        VALUES (?, ?, ?, ?, NOW())
        ON CONFLICT (store_id, product_id)
        DO UPDATE SET quantity = store_inventories.quantity + ?, updated_by = ?, updated_at = NOW()
    `, storeID, productID, delta, updatedBy, delta, updatedBy).Error
}

func (r *storeInventoryRepository) GetAllInventory() ([]*domain.StoreInventory, error) {
	var inventories []*domain.StoreInventory
	err := r.db.Preload("Store").Preload("Product").Find(&inventories).Error
	return inventories, err
}
