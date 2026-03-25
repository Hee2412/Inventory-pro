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

type AdminOrderService interface {
	GetAllOrderInSession(sessionID uint) ([]*response.AdminOrderInSessionResponse, error)
	ApproveOrder(orderId uint) error
	DeclineOrder(orderId uint, reason string) error
	GetAllPaginatedSessions(params request.OrderSearchParams) ([]*response.OrderResponse, int64, error)
	GetStoreWithOutOrder(sessionID uint) (*response.StoreWithoutOrderResponse, error)
	DeliverOrder(orderId uint) error
	RedeliverOrder(orderId uint, req request.RedeliverRequest) error
}

type adminOrderService struct {
	orderSessionRepo   repository.OrderSessionRepository
	storeOrderRepo     repository.StoreOrderRepository
	userRepo           repository.UserRepository
	storeOrderItemRepo repository.StoreOrderItemRepository
}

func NewAdminOrderService(
	orderSessionRepo repository.OrderSessionRepository,
	storeOrderRepo repository.StoreOrderRepository,
	userRepo repository.UserRepository,
	storeOrderItemRepo repository.StoreOrderItemRepository) AdminOrderService {
	return &adminOrderService{
		orderSessionRepo:   orderSessionRepo,
		storeOrderRepo:     storeOrderRepo,
		userRepo:           userRepo,
		storeOrderItemRepo: storeOrderItemRepo,
	}
}

func (a *adminOrderService) GetAllOrderInSession(sessionID uint) ([]*response.AdminOrderInSessionResponse, error) {
	_, err := a.orderSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	orders, err := a.storeOrderRepo.FindBySessionID(sessionID)
	if err != nil {
		return make([]*response.AdminOrderInSessionResponse, 0), nil
	}
	result := make([]*response.AdminOrderInSessionResponse, 0)
	for _, order := range orders {
		result = append(result, toOrderInSessionResponse(order))
	}
	return result, nil
}

func (a *adminOrderService) ApproveOrder(orderId uint) error {
	order, err := a.storeOrderRepo.FindById(orderId)
	if err != nil {
		return domain.ErrNotFound
	}
	if order.Status != "SUBMITTED" && order.Status != "NO_ORDER" {
		return fmt.Errorf("%w: cannot approve order with current status: %s", domain.ErrInvalidInput, order.Status)
	}
	if order.Status == "UNSUBMITTED_EXPIRED" {
		return fmt.Errorf("%w: this order has expired without store interaction", domain.ErrInvalidInput)
	}
	order.Status = "APPROVED"
	now := time.Now()
	order.ApproveAt = &now
	err = a.storeOrderRepo.Update(order)
	if err != nil {
		return fmt.Errorf("%w: failed to approve order: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (a *adminOrderService) DeclineOrder(orderId uint, reason string) error {
	order, err := a.storeOrderRepo.FindById(orderId)
	if err != nil {
		return domain.ErrNotFound
	}
	if order.Status == "DRAFT" {
		return fmt.Errorf("%w: only reports in SUBMITTED status can be decline", domain.ErrInvalidInput)
	}
	order.Status = "DECLINED"
	order.Note = reason
	err = a.storeOrderRepo.Update(order)
	if err != nil {
		return fmt.Errorf("%w: failed to update status: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (a *adminOrderService) GetAllPaginatedSessions(params request.OrderSearchParams) ([]*response.OrderResponse, int64, error) {
	sessions, total, err := a.storeOrderRepo.SearchAndFilter(params)
	if err != nil {
		return nil, 0, err
	}
	result := make([]*response.OrderResponse, 0, len(sessions))
	for _, session := range sessions {
		result = append(result, toOrderResponse(session))
	}
	return result, total, nil
}

func (a *adminOrderService) GetStoreWithOutOrder(sessionID uint) (*response.StoreWithoutOrderResponse, error) {
	stores, err := a.userRepo.FindByRoleAndActive("store", true)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to fetch data: %v", domain.ErrDatabase, err)
	}
	session, err := a.orderSessionRepo.FindById(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrSessionNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch session data: %v", domain.ErrDatabase, err)
	}
	orders, err := a.storeOrderRepo.FindBySessionID(sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, fmt.Errorf("%w: failed to fetch order: %v", domain.ErrDatabase, err)
	}
	confirmedOrder := make(map[uint]bool)
	for _, order := range orders {
		if order.ConfirmedAt != nil {
			confirmedOrder[order.StoreID] = true
		}
	}
	var notConfirmedStores []*response.StoreTrackingResponse
	for _, store := range stores {
		if !confirmedOrder[store.ID] {
			notConfirmedStores = append(notConfirmedStores, &response.StoreTrackingResponse{
				StoreID:   store.ID,
				StoreName: store.StoreName,
			})
		}
	}
	return &response.StoreWithoutOrderResponse{
		SessionID:         sessionID,
		SessionName:       session.Title,
		NotOrdered:        len(notConfirmedStores),
		StoreWithoutOrder: notConfirmedStores,
	}, nil
}

func (a *adminOrderService) DeliverOrder(orderId uint) error {
	order, err := a.storeOrderRepo.FindById(orderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%w: failed to fetch order: %v", domain.ErrDatabase, orderId)
	}
	if order.Status != "APPROVED" {
		return fmt.Errorf("%w: only APPROVED orders can be delivered, current: %s",
			domain.ErrInvalidInput, order.Status)
	}
	now := time.Now()
	order.Status = "DELIVERED"
	order.ConfirmedAt = &now
	err = a.storeOrderRepo.Update(order)
	if err != nil {
		return fmt.Errorf("%w: failed to deliver order: %v", domain.ErrDatabase, err)
	}
	return nil
}

func (a *adminOrderService) RedeliverOrder(orderId uint, req request.RedeliverRequest) error {
	order, err := a.storeOrderRepo.FindById(orderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.ErrNotFound
		}
		return fmt.Errorf("%w: failed to fetch order: %v", domain.ErrDatabase, orderId)
	}
	if order.Status != "REJECTED" {
		return fmt.Errorf("%w: only REJECTED orders can be redelivered, current: %s",
			domain.ErrInvalidInput, order.Status)
	}

	if err = a.storeOrderItemRepo.DeleteByOrderId(orderId); err != nil {
		return fmt.Errorf("%w: failed to clear old items: %v", domain.ErrDatabase, err)
	}

	var newItems []*domain.OrderItems
	for _, item := range req.Items {
		if item.Quantity == 0 {
			continue
		}
		newItems = append(newItems, &domain.OrderItems{
			OrderID:   orderId,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	if len(newItems) == 0 {
		return fmt.Errorf("%w: redeliver must have at least one item with quantity > 0", domain.ErrInvalidInput)
	}
	if err = a.storeOrderItemRepo.Create(newItems...); err != nil {
		return fmt.Errorf("%w: failed to create new items: %v", domain.ErrDatabase, err)
	}

	now := time.Now()
	order.Status = "DELIVERED"
	order.ConfirmedAt = &now
	order.Note = ""
	if err = a.storeOrderRepo.Update(order); err != nil {
		return fmt.Errorf("%w: failed to update order status: %v", domain.ErrDatabase, err)
	}
	return nil
}
