package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"fmt"
	"time"
)

type AdminOrderService interface {
	GetAllOrderInSession(sessionID uint) ([]*response.AdminOrderInSessionResponse, error)
	ApproveOrder(orderId uint) error
	DeclineOrder(orderId uint, reason string) error
	GetAllPaginatedSessions(params request.OrderSearchParams) ([]*response.OrderResponse, int64, error)
	GetStoreWithOutOrder(sessionID uint) (*response.StoreWithoutOrderResponse, error)
}

type adminOrderService struct {
	orderSessionRepo repository.OrderSessionRepository
	storeOrderRepo   repository.StoreOrderRepository
	userRepo         repository.UserRepository
}

func NewAdminOrderService(
	orderSessionRepo repository.OrderSessionRepository,
	storeOrderRepo repository.StoreOrderRepository,
	userRepo repository.UserRepository) AdminOrderService {
	return &adminOrderService{
		orderSessionRepo: orderSessionRepo,
		storeOrderRepo:   storeOrderRepo,
		userRepo:         userRepo}
}

func toOrderInSessionResponse(order *domain.StoreOrder) *response.AdminOrderInSessionResponse {
	return &response.AdminOrderInSessionResponse{
		ID:          order.ID,
		StoreID:     order.StoreID,
		StoreName:   order.StoreName,
		Status:      order.Status,
		SubmittedAt: order.SubmittedAt,
		CreatedAt:   order.CreatedAt,
	}
}

func toOrderResponse(order *domain.StoreOrder) *response.OrderResponse {
	return &response.OrderResponse{
		ID:        order.ID,
		StoreID:   order.StoreID,
		StoreName: order.StoreName,
		Status:    order.Status,
		SessionID: order.SessionID,
	}
}

func (a *adminOrderService) GetAllOrderInSession(sessionID uint) ([]*response.AdminOrderInSessionResponse, error) {
	//var sessionID
	_, err := a.orderSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	//mapping to orderSessionResponse
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
	//Find orderByID
	order, err := a.storeOrderRepo.FindById(orderId)
	if err != nil {
		return domain.ErrNotFound
	}
	//Check status ("Submitted")
	if order.Status != "SUBMITTED" && order.Status != "NO_ORDER" {
		return fmt.Errorf("%w: cannot approve order with current status: %s", domain.ErrInvalidInput, order.Status)
	}
	if order.Status == "UNSUBMITTED_EXPIRED" {
		return fmt.Errorf("%w: this order has expired without store interaction", domain.ErrInvalidInput)
	}
	//Change status ("Approved")
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
	//Find orderByID
	order, err := a.storeOrderRepo.FindById(orderId)
	if err != nil {
		return domain.ErrNotFound
	}
	//Check status ("Submitted")
	if order.Status == "DRAFT" {
		return fmt.Errorf("%w: only reports in SUBMITTED status can be decline", domain.ErrInvalidInput)
	}
	//Change status ("Declined")
	order.Status = "DECLINED"
	order.Note = reason
	//Save/Update
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
		return nil, domain.ErrDatabase
	}
	session, err := a.orderSessionRepo.FindById(sessionID)
	if err != nil {
		return nil, domain.ErrNotFound
	}
	orders, err := a.storeOrderRepo.FindBySessionID(sessionID)
	if err != nil {
		return nil, domain.ErrNotFound
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
