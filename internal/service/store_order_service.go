package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"time"
)

type StoreOrderService interface {
	GetOrCreateOrder(sessionID uint, storeID uint) (*response.StoreOrderDetailResponse, error)
	UpdateOrder(orderID uint, req request.UpdateOrderItemRequest) error
	SubmitOrder(orderId uint) error
	GetOrderDetail(orderID uint) (*response.StoreOrderDetailResponse, error)
	GetMyOrder(storeID uint) ([]response.StoreOrderResponse, error)
}

type storeOrderService struct {
	storeOrderRepo     repository.StoreOrderRepository
	orderSessionRepo   repository.OrderSessionRepository
	storeOrderItemRepo repository.StoreOrderItemRepository
}

func NewStoreOrderService(
	storeOrderRepo repository.StoreOrderRepository,
	orderSessionRepo repository.OrderSessionRepository,
	storeOrderItemRepo repository.StoreOrderItemRepository,
) StoreOrderService {
	return &storeOrderService{
		storeOrderRepo:     storeOrderRepo,
		orderSessionRepo:   orderSessionRepo,
		storeOrderItemRepo: storeOrderItemRepo,
	}
}

func toOrderItemResponse(items []domain.OrderItems) []response.OrderItemResponse {
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

func toStoreOrderResponse(order *domain.StoreOrder) response.StoreOrderResponse {
	return response.StoreOrderResponse{
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
	session, err := s.orderSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}
	if session.Status != "OPEN" {
		return nil, errors.New("session is not open")
	}
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
	order, err := s.storeOrderRepo.FindById(orderID)
	if err != nil {
		return errors.New("order not found")
	}
	if order.Status != "DRAFT" {
		return errors.New("can only edit draft orders")
	}
	items, err := s.storeOrderItemRepo.FindByOrderId(orderID)
	if err != nil {
		return err
	}
	var targetItem *domain.OrderItems
	for i, item := range items {
		if item.ProductID == req.ProductID {
			targetItem = &items[i]
			break
		}
	}
	if targetItem == nil {
		newItem := &domain.OrderItems{
			OrderID:   order.ID,
			Quantity:  req.Quantity,
			ProductID: req.ProductID}
		return s.storeOrderItemRepo.Create(newItem)
	}
	targetItem.Quantity = req.Quantity
	return s.storeOrderItemRepo.Update(targetItem)
}

func (s *storeOrderService) SubmitOrder(orderId uint) error {
	order, err := s.storeOrderRepo.FindById(orderId)
	if err != nil {
		return errors.New("order not found")
	}
	if order.Status != "DRAFT" {
		return errors.New("can only edit draft orders")
	}
	items, err := s.storeOrderItemRepo.FindByOrderId(orderId)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return errors.New("cannot submit empty order")
	}
	now := time.Now()
	order.Status = "SUBMITTED"
	order.SubmittedAt = &now
	return s.storeOrderRepo.Update(order)
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

func (s *storeOrderService) GetMyOrder(storeID uint) ([]response.StoreOrderResponse, error) {
	//findOrder By storeID
	orders, err := s.storeOrderRepo.FindByStoreID(storeID)
	if err != nil {
		return make([]response.StoreOrderResponse, 0), nil
	}
	//transfer to storeOrderResponse
	result := make([]response.StoreOrderResponse, 0)
	for _, order := range orders {
		result = append(result, toStoreOrderResponse(order))
	}
	//return result
	return result, nil
}
