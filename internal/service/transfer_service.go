package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type TransferService interface {
	CreateTransfer(req request.CreateTransferRequest, createdBy uint, fromStoreID uint) (*response.TransferOrderDetailResponse, error)
	GetTransferDetail(transferID uint) (*response.TransferOrderDetailResponse, error)
	GetMyTransfers(storeID uint, direction string) ([]*response.TransferOrderResponse, error)
	GetAllTransfers(page, limit int, status string) ([]*response.TransferOrderResponse, int64, error)
	ApproveTransfer(transferID uint, approvedBy uint) (*response.TransferOrderDetailResponse, error)
	CancelTransfer(transferID uint, cancelledBy uint, reason string) error
}

type transferService struct {
	transferRepo     repository.TransferOrderRepository
	transferItemRepo repository.TransferOrderItemRepository
	inventoryRepo    repository.StoreInventoryRepository
	userRepo         repository.UserRepository
	productRepo      repository.ProductRepository
}

func NewTransferService(
	transferRepo repository.TransferOrderRepository,
	transferItemRepo repository.TransferOrderItemRepository,
	inventoryRepo repository.StoreInventoryRepository,
	userRepo repository.UserRepository,
	productRepo repository.ProductRepository,
) TransferService {
	return &transferService{
		transferRepo:     transferRepo,
		transferItemRepo: transferItemRepo,
		inventoryRepo:    inventoryRepo,
		userRepo:         userRepo,
		productRepo:      productRepo,
	}
}

// CreateTransfer
// Rule: Store create transfer from self - admin can create transfer from any - storeID take from req if admin create
func (s *transferService) CreateTransfer(req request.CreateTransferRequest, createdBy uint, fromStoreID uint) (*response.TransferOrderDetailResponse, error) {
	//validate store cant not transfer to self
	if fromStoreID == req.ToStoreID {
		return nil, fmt.Errorf("%w: cannot transfer to the same store", domain.ErrInvalidInput)
	}
	_, err := s.userRepo.FindById(req.ToStoreID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w: destination store not found", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("%w: failed to validate store: %v", domain.ErrDatabase, err)
	}

	//create transfer order
	order := &domain.TransferOrder{
		FromStoreID: fromStoreID,
		ToStoreID:   req.ToStoreID,
		Status:      "PENDING",
		Note:        req.Note,
		CreatedBy:   createdBy,
	}
	if err = s.transferRepo.Create(order); err != nil {
		return nil, fmt.Errorf("%w: failed to create transfer: %v", domain.ErrDatabase, err)
	}

	//create items
	var items []*domain.TransferOrderItem
	for _, item := range req.Items {
		product, err := s.productRepo.FindById(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("%w: product id %d not found", domain.ErrNotFound, item.ProductID)
		}
		items = append(items, &domain.TransferOrderItem{
			TransferOrderID: order.ID,
			ProductID:       product.ID,
			ProductName:     product.ProductName,
			ProductCode:     product.ProductCode,
			Quantity:        item.Quantity,
		})
	}
	if err = s.transferItemRepo.Create(items...); err != nil {
		return nil, fmt.Errorf("%w: failed to create transfer items: %v", domain.ErrDatabase, err)
	}

	return s.buildDetailResponse(order, items)
}

// ApproveTransfer
// Store accept transfer order or admin/ minus inventory fromStore and plus inventory toStore
func (s *transferService) ApproveTransfer(transferID uint, approvedBy uint) (*response.TransferOrderDetailResponse, error) {
	order, err := s.transferRepo.FindById(transferID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch transfer: %v", domain.ErrDatabase, transferID)
	}
	if order.Status != "PENDING" {
		return nil, fmt.Errorf("%w: only PENDING transfers can be approved, current: %s",
			domain.ErrInvalidInput, order.Status)
	}

	items, err := s.transferItemRepo.FindByTransferOrderId(transferID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch transfer items: %v", domain.ErrDatabase, transferID)
	}

	// Update inventory for each item
	for _, item := range items {
		// Minus FromStore (delta -)
		if err = s.inventoryRepo.AdjustQuantity(order.FromStoreID, item.ProductID, -item.Quantity, approvedBy); err != nil {
			return nil, fmt.Errorf("%w: failed to deduct inventory from store %d for product %d: %v",
				domain.ErrDatabase, order.FromStoreID, item.ProductID, err)
		}
		// Plus ToStore (delta +)
		if err = s.inventoryRepo.AdjustQuantity(order.ToStoreID, item.ProductID, item.Quantity, approvedBy); err != nil {
			return nil, fmt.Errorf("%w: failed to add inventory to store %d for product %d: %v",
				domain.ErrDatabase, order.ToStoreID, item.ProductID, err)
		}
	}

	now := time.Now()
	order.Status = "APPROVED"
	order.ApprovedAt = &now
	order.ApprovedBy = &approvedBy
	if err = s.transferRepo.Update(order); err != nil {
		return nil, fmt.Errorf("%w: failed to update transfer status: %v", domain.ErrDatabase, err)
	}

	return s.buildDetailResponse(order, items)
}

// CancelTransfer
func (s *transferService) CancelTransfer(transferID uint, cancelledBy uint, reason string) error {
	order, err := s.transferRepo.FindById(transferID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%w: failed to fetch transfer: %v", domain.ErrDatabase, transferID)
	}
	if order.Status != "PENDING" {
		return fmt.Errorf("%w: only PENDING transfers can be cancelled, current: %s",
			domain.ErrInvalidInput, order.Status)
	}
	now := time.Now()
	order.Status = "CANCELLED"
	order.CancelledAt = &now
	order.CancelReason = reason
	if err = s.transferRepo.Update(order); err != nil {
		return fmt.Errorf("%w: failed to cancel transfer: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (s *transferService) GetTransferDetail(transferID uint) (*response.TransferOrderDetailResponse, error) {
	order, err := s.transferRepo.FindById(transferID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch transfer: %v", domain.ErrDatabase, transferID)
	}
	items, err := s.transferItemRepo.FindByTransferOrderId(transferID)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch transfer items: %v", domain.ErrDatabase, transferID)
	}
	return s.buildDetailResponse(order, items)
}

func (s *transferService) GetMyTransfers(storeID uint, direction string) ([]*response.TransferOrderResponse, error) {
	var orders []*domain.TransferOrder
	var err error
	switch direction {
	case "out":
		orders, err = s.transferRepo.FindByFromStore(storeID)
	case "in":
		orders, err = s.transferRepo.FindByToStore(storeID)
	default:
		outOrders, errOut := s.transferRepo.FindByFromStore(storeID)
		inOrders, errIn := s.transferRepo.FindByToStore(storeID)
		if errOut != nil || errIn != nil {
			return make([]*response.TransferOrderResponse, 0), nil
		}
		orders = append(outOrders, inOrders...)
	}
	if err != nil {
		return make([]*response.TransferOrderResponse, 0), nil
	}
	result := make([]*response.TransferOrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, toTransferOrderResponse(o))
	}
	return result, nil
}

func (s *transferService) GetAllTransfers(page, limit int, status string) ([]*response.TransferOrderResponse, int64, error) {
	orders, total, err := s.transferRepo.FindAllPaginated(page, limit, status)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: failed to fetch transfers: %v", domain.ErrDatabase, err)
	}
	result := make([]*response.TransferOrderResponse, 0, len(orders))
	for _, o := range orders {
		result = append(result, toTransferOrderResponse(o))
	}
	return result, total, nil
}

func (s *transferService) buildDetailResponse(order *domain.TransferOrder, items []*domain.TransferOrderItem) (*response.TransferOrderDetailResponse, error) {
	return &response.TransferOrderDetailResponse{
		Order: toTransferOrderResponse(order),
		Items: toTransferItemResponse(items),
	}, nil
}
