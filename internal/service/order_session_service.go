package service

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
)

type OrderSessionService interface {
	CreateSession(req request.CreateOrderSessionRequest) (*response.OrderSessionResponse, error)
	GetAllSessions() ([]response.OrderSessionResponse, error)
	GetSessionById(sessionId uint) (*response.StoreOrderDetailResponse, error)
	AddProductToSession(req request.AddProductToSessionRequest) error
	RemoveProductFromSession(req request.UpdateOrderItemRequest) error
	CloseSession(sessionId uint) error
}
