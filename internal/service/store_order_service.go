package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"fmt"
	"math"
	"time"
)

type StoreOrderService interface {
	GetOrCreateOrder(sessionID uint, storeID uint) (*response.StoreOrderDetailResponse, error)
	UpdateOrder(orderID uint, req request.UpdateOrderItemRequest) error
	GetOrderDetail(orderID uint) (*response.StoreOrderDetailResponse, error)
	GetMyOrder(storeID uint) ([]*response.StoreOrderResponse, error)
	GetAllPaginatedOrders(params request.OrderSearchParams) ([]*response.StoreOrderResponse, int64, error)
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

func toOrderItemResponse(items []*domain.OrderItems) []response.OrderItemResponse {
	result := make([]response.OrderItemResponse, 0)
	for _, item := range items {
		result = append(result, response.OrderItemResponse{
			ID:          item.ID,
			OrderID:     item.OrderID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}
	return result
}

func toStoreOrderResponse(order *domain.StoreOrder) *response.StoreOrderResponse {
	return &response.StoreOrderResponse{
		ID:          order.ID,
		SessionID:   order.SessionID,
		StoreID:     order.StoreID,
		Status:      order.Status,
		SubmittedAt: order.SubmittedAt,
		ApprovedAt:  order.ApproveAt,
		CreatedAt:   order.CreatedAt,
	}
}

func (s *storeOrderService) GetOrCreateOrder(sessionID uint, storeID uint) (*response.StoreOrderDetailResponse, error) {
	//call session
	session, err := s.orderSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}
	//check deadline
	if session.Status == "OPEN" && time.Now().After(session.Deadline) {
		session.Status = "CLOSED"
		err := s.orderSessionRepo.Update(session)
		if err != nil {
			return nil, err
		}
		return nil, errors.New("session is expired and been closed")
	}
	//check session status
	if session.Status == "CLOSED" {
		return nil, errors.New("session is closed")
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
			return nil, err
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
		return errors.New("order not found")
	}
	//check session
	session, err := s.orderSessionRepo.FindById(order.SessionID)
	if err != nil {
		return errors.New("session not found")
	}
	if session.Status == "OPEN" && time.Now().After(session.Deadline) {
		session.Status = "CLOSED"
		err := s.orderSessionRepo.Update(session)
		if err != nil {
			return err
		}
		return errors.New("session expired and been closed")
	}
	if session.Status == "CLOSED" {
		return errors.New("session is already closed")
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
		return errors.New("no valid items found")
	}
	err = s.storeOrderItemRepo.DeleteByOrderId(orderID)
	if err != nil {
		return err
	}
	return s.storeOrderItemRepo.Create(newItems...)
}

func (s *storeOrderService) GetOrderDetail(OrderID uint) (*response.StoreOrderDetailResponse, error) {
	//Find Order
	order, err := s.storeOrderRepo.FindById(OrderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	//Check item
	items, err := s.storeOrderItemRepo.FindByOrderId(OrderID)
	if err != nil {
		return nil, err
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
