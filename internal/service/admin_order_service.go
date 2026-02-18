package service

import "Inventory-pro/internal/dto/response"

type AdminOrderService interface {
	GetAllOrderInSession(sessionID uint) (*[]response.OrderSessionResponse, error)
	ApproveOrder(orderId uint) error
	DeclineOrder(orderId uint, reason string) error
}
