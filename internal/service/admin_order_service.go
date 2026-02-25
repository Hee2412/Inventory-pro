package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/repository"
	"errors"
	"time"
)

type AdminOrderService interface {
	GetAllOrderInSession(sessionID uint) ([]response.AdminOrderInSessionResponse, error)
	ApproveOrder(orderId uint) error
	DeclineOrder(orderId uint, reason string) error
}

type adminOrderService struct {
	sessionRepo repository.OrderSessionRepository
	orderRepo   repository.StoreOrderRepository
}

func NewAdminOrderService(
	repo repository.OrderSessionRepository,
	orderRepo repository.StoreOrderRepository) AdminOrderService {
	return &adminOrderService{
		sessionRepo: repo,
		orderRepo:   orderRepo}
}

func toOrderInSessionResponse(order *domain.StoreOrder) response.AdminOrderInSessionResponse {
	return response.AdminOrderInSessionResponse{
		ID:          order.ID,
		StoreID:     order.StoreID,
		StoreName:   order.StoreName,
		Status:      order.Status,
		SubmittedAt: order.SubmittedAt,
		CreatedAt:   order.CreatedAt,
	}
}

func (a *adminOrderService) GetAllOrderInSession(sessionID uint) ([]response.AdminOrderInSessionResponse, error) {
	//var sessionID
	_, err := a.sessionRepo.FindById(sessionID)
	if err != nil {
		return nil, errors.New("session not found")
	}
	//mapping to orderSessionResponse
	orders, err := a.orderRepo.FindBySessionID(sessionID)
	if err != nil {
		return make([]response.AdminOrderInSessionResponse, 0), nil
	}
	result := make([]response.AdminOrderInSessionResponse, 0)
	for _, order := range orders {
		result = append(result, toOrderInSessionResponse(order))
	}
	return result, nil
}

func (a *adminOrderService) ApproveOrder(orderId uint) error {
	//Find orderByID
	order, err := a.orderRepo.FindById(orderId)
	if err != nil {
		return errors.New("order not found")
	}
	//Check status ("Submitted")
	if order.Status != "SUBMITTED" {
		return errors.New("order status is not SUBMITTED")
	}
	//Change status ("Approved")
	order.Status = "APPROVED"
	//Approve By/At
	var now = time.Now()
	order.ApproveAt = &now
	//Save/Update
	return a.orderRepo.Update(order)
}

func (a *adminOrderService) DeclineOrder(orderId uint, reason string) error {
	//Find orderByID
	order, err := a.orderRepo.FindById(orderId)
	if err != nil {
		return errors.New("order not found")
	}
	//Check status ("Submitted")
	if order.Status != "SUBMITTED" {
		return errors.New("order status is not SUBMITTED")
	}
	//Change status ("Declined")
	order.Status = "DECLINED"
	order.Note = reason
	//Save/Update
	return a.orderRepo.Update(order)
}
