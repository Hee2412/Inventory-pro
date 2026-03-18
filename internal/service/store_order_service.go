package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"math"
	"time"
)

type StoreOrderService interface {
	GetOrCreateOrder(sessionID uint, storeID uint) (*response.StoreOrderDetailResponse, error)
	UpdateOrder(orderID uint, req request.UpdateOrderItemRequest) error
	GetOrderDetail(orderID uint) (*response.StoreOrderDetailResponse, error)
	GetMyOrder(storeID uint) ([]*response.StoreOrderResponse, error)
	GetAllPaginatedOrders(params request.OrderSearchParams) ([]*response.StoreOrderResponse, int64, error)
	UpdateStatus(OrderID uint) (*domain.StoreOrder, error)
}

type storeOrderService struct {
	storeOrderRepo     repository.StoreOrderRepository
	orderSessionRepo   repository.OrderSessionRepository
	storeOrderItemRepo repository.StoreOrderItemRepository
	productRepo        repository.ProductRepository
}

func NewStoreOrderService(
	storeOrderRepo repository.StoreOrderRepository,
	orderSessionRepo repository.OrderSessionRepository,
	storeOrderItemRepo repository.StoreOrderItemRepository,
	productRepo repository.ProductRepository,
) StoreOrderService {
	return &storeOrderService{
		storeOrderRepo:     storeOrderRepo,
		orderSessionRepo:   orderSessionRepo,
		storeOrderItemRepo: storeOrderItemRepo,
		productRepo:        productRepo,
	}
}

func (s *storeOrderService) GetOrCreateOrder(sessionID uint, storeID uint) (*response.StoreOrderDetailResponse, error) {
	//call session
	session, err := s.orderSessionRepo.FindById(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session: %v", domain.ErrDatabase, session.ID)
	}
	//check session status
	if session.Status == "CLOSED" {
		return nil, fmt.Errorf("%w: session is closed", domain.ErrSessionClosed)
	}
	//check deadline
	if time.Now().After(session.Deadline) {
		order, err := s.storeOrderRepo.FindByStoreAndSession(sessionID, storeID)
		if err != nil {
			order = &domain.StoreOrder{
				SessionID: sessionID,
				StoreID:   storeID,
				Status:    "UNSUBMITTED_EXPIRED",
			}
			_ = s.storeOrderRepo.Create(order)
		} else if order.Status == "DRAFT" {
			order.Status = "UNSUBMITTED_EXPIRED"
			_ = s.storeOrderRepo.Update(order)
		}
		return nil, fmt.Errorf("%w: session is closed", domain.ErrSessionClosed)
	}

	//check condition get or create order
	order, err := s.storeOrderRepo.FindByStoreAndSession(sessionID, storeID)
	if err != nil {
		order = &domain.StoreOrder{
			SessionID: sessionID,
			StoreID:   storeID,
			Status:    "DRAFT",
		}
		err = s.storeOrderRepo.Create(order)
		if err != nil {
			return nil, fmt.Errorf("%w: failed to store order: %v", domain.ErrDatabase, err)
		}
	}
	items, _ := s.storeOrderItemRepo.FindByOrderId(order.ID)
	result := &response.StoreOrderDetailResponse{
		Order: toStoreOrderResponse(order),
		Items: toOrderItemResponse(items),
	}
	return result, nil
}

func (s *storeOrderService) UpdateOrder(orderID uint, req request.UpdateOrderItemRequest) error {
	//find order
	order, err := s.storeOrderRepo.FindById(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%w: failed to fetch order: %v", domain.ErrDatabase, orderID)
	}
	//check session
	session, err := s.orderSessionRepo.FindById(order.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrSessionNotFound
		}
		return fmt.Errorf("%w: failed to fetch session: %v", domain.ErrDatabase, session.ID)
	}
	if session.Status == "CLOSED" {
		return fmt.Errorf("%w: session is closed", domain.ErrSessionClosed)
	}
	if time.Now().After(session.Deadline) {
		if order.Status == "DRAFT" {
			order.Status = "UNSUBMITTED_EXPIRED"
			_ = s.storeOrderRepo.Update(order)
		}
		return fmt.Errorf("%w: deadline exceeded at %v", domain.ErrSessionClosed, session.Deadline)
	}

	//var order items
	var newItems []*domain.OrderItems
	for _, item := range req.Items {
		if item.Quantity == 0 {
			continue
		}
		product, err := s.productRepo.FindById(item.ProductID)
		if err != nil {
			continue
		}
		//check moq
		if item.Quantity < product.MOQ {
			return fmt.Errorf("product %s: quantity %.0f is less than MOQ %.0f",
				product.ProductName, item.Quantity, product.MOQ)
		}
		//check om
		if product.OM > 0 {
			if math.Mod(item.Quantity, product.OM) != 0 {
				return fmt.Errorf("product %s: quantity must be multiple of %.0f (e.g., %.0f, %.0f, %.0f...)",
					product.ProductName, product.OM, product.OM, product.OM*2, product.OM*3)
			}
		}
		newItems = append(newItems, &domain.OrderItems{
			OrderID:     order.ID,
			ProductID:   product.ID,
			Quantity:    item.Quantity,
			ProductName: product.ProductName,
			ProductCode: product.ProductCode,
		})
	}
	if len(newItems) == 0 {
		return domain.ErrInvalidInput
	}
	err = s.storeOrderItemRepo.DeleteByOrderId(orderID)
	if err != nil {
		return fmt.Errorf("%w: clear old items failed: %v", domain.ErrDatabase, err)
	}
	now := time.Now()
	order.Status = "SUBMITTED"
	order.ConfirmedAt = &now
	err = s.storeOrderRepo.Update(order)
	if err != nil {
		return fmt.Errorf("%w: update failed: %v", domain.ErrDatabase, err)
	}
	err = s.storeOrderItemRepo.Create(newItems...)
	if err != nil {
		return fmt.Errorf("%w: create new items failed: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (s *storeOrderService) GetOrderDetail(orderID uint) (*response.StoreOrderDetailResponse, error) {
	//Find Order
	order, err := s.storeOrderRepo.FindById(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch order: %v", domain.ErrDatabase, orderID)
	}
	//Check item
	items, err := s.storeOrderItemRepo.FindByOrderId(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch items in order: %v", domain.ErrDatabase, orderID)
	}
	//transfer orderItem to response
	result := &response.StoreOrderDetailResponse{
		Order: toStoreOrderResponse(order),
		Items: toOrderItemResponse(items),
	}
	//return result
	return result, nil
}

func (s *storeOrderService) GetMyOrder(storeID uint) ([]*response.StoreOrderResponse, error) {
	//findOrder By storeID
	orders, err := s.storeOrderRepo.FindByStoreID(storeID)
	if err != nil {
		return make([]*response.StoreOrderResponse, 0), nil
	}
	//transfer to storeOrderResponse
	result := make([]*response.StoreOrderResponse, 0)
	for _, order := range orders {
		result = append(result, toStoreOrderResponse(order))
	}
	//return result
	return result, nil
}

func (s *storeOrderService) GetAllPaginatedOrders(params request.OrderSearchParams) ([]*response.StoreOrderResponse, int64, error) {
	orders, total, err := s.storeOrderRepo.SearchAndFilter(params)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*response.StoreOrderResponse, 0)
	for _, order := range orders {
		result = append(result, toStoreOrderResponse(order))
	}
	return result, total, nil
}

func (s *storeOrderService) UpdateStatus(orderID uint) (*domain.StoreOrder, error) {
	order, err := s.storeOrderRepo.FindById(orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch order: %v", domain.ErrDatabase, orderID)
	}
	session, err := s.orderSessionRepo.FindById(order.SessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session: %v", domain.ErrDatabase, session.ID)
	}
	if session.Status == "CLOSE" {
		return nil, fmt.Errorf("%w: session is closed", domain.ErrSessionClosed)
	}
	now := time.Now()
	if order.Status == "DRAFT" {
		order.Status = "NO_ORDER"
		order.ConfirmedAt = &now
	} else if order.Status == "NO_ORDER" {
		order.Status = "DRAFT"
		order.ConfirmedAt = nil
		order.UpdatedAt = &now
	}
	err = s.storeOrderRepo.Update(order)
	if err != nil {
		return nil, fmt.Errorf("%w: update failed: %v", domain.ErrDatabase, err)
	}
	return order, nil
}
