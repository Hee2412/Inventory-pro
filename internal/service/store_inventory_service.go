package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type InventoryService interface {
	GetStoreInventory(storeID uint) ([]*response.InventoryResponse, error)
	GetProductInventoryAcrossStores(productID uint) ([]*response.InventoryResponse, error)
	UpdateInventory(storeID uint, req request.UpdateInventoryRequest, updatedBy uint) error
	BatchAdjustInventory(req []request.AdjustmentItem, updatedBy uint) error
	GetAllInventory() ([]*response.InventoryResponse, error)
}

type inventoryService struct {
	inventoryRepo repository.StoreInventoryRepository
	storeRepo     repository.UserRepository
	productRepo   repository.ProductRepository
}

func NewInventoryService(
	inventoryRepo repository.StoreInventoryRepository,
	storeRepo repository.UserRepository,
	productRepo repository.ProductRepository,
) InventoryService {
	return &inventoryService{
		inventoryRepo: inventoryRepo,
		storeRepo:     storeRepo,
		productRepo:   productRepo,
	}
}

func (s *inventoryService) GetProductInventoryAcrossStores(productID uint) ([]*response.InventoryResponse, error) {
	_, err := s.productRepo.FindById(productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrProductNotFound
		}
		return nil, fmt.Errorf("%w: failed to find product: %v", domain.ErrDatabase, err)
	}

	inventories, err := s.inventoryRepo.FindByProductID(productID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to get inventory: %v", domain.ErrDatabase, err)
	}

	result := make([]*response.InventoryResponse, 0, len(inventories))
	for _, inv := range inventories {
		result = append(result, toInventoryResponse(inv))
	}

	return result, nil
}

func (s *inventoryService) GetAllInventory() ([]*response.InventoryResponse, error) {
	inventories, err := s.inventoryRepo.GetAllInventory()
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch data", domain.ErrDatabase)
	}
	result := make([]*response.InventoryResponse, 0, len(inventories))
	for _, product := range inventories {
		result = append(result, toInventoryResponse(product))
	}
	return result, nil
}

func (s *inventoryService) GetStoreInventory(storeID uint) ([]*response.InventoryResponse, error) {
	inventories, err := s.inventoryRepo.FindByStoreID(storeID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrDatabase, err)
	}

	result := make([]*response.InventoryResponse, 0, len(inventories))
	for _, inv := range inventories {
		result = append(result, toInventoryResponse(inv))
	}
	return result, nil
}

func (s *inventoryService) UpdateInventory(storeID uint, req request.UpdateInventoryRequest, updatedBy uint) error {
	// Validate store exists
	_, err := s.storeRepo.FindById(storeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrUserNotFound
		}
		return fmt.Errorf("%w: failed to find store: %v", domain.ErrDatabase, err)
	}

	// Validate product exists
	_, err = s.productRepo.FindById(req.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrProductNotFound
		}
		return fmt.Errorf("%w: productId invalid %v", domain.ErrDatabase, err)
	}
	if req.Quantity < 0 {
		return fmt.Errorf("%w: quantity cannot be negative <0", domain.ErrInvalidInput)
	}

	// Update or create inventory
	inventory := &domain.StoreInventory{
		StoreID:   storeID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		UpdatedBy: updatedBy,
	}

	err = s.inventoryRepo.Upsert(inventory)
	if err != nil {
		return fmt.Errorf("%w: failed to update store: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (s *inventoryService) BatchAdjustInventory(req []request.AdjustmentItem, updatedBy uint) error {
	for _, adj := range req {
		err := s.inventoryRepo.AdjustQuantity(adj.StoreID, adj.ProductID, adj.Delta, updatedBy)
		if err != nil {
			return fmt.Errorf("%w: failed to adjust inventory for store %d product %d: %v",
				domain.ErrDatabase, adj.StoreID, adj.ProductID, err)
		}
	}
	return nil
}
