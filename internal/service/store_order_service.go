package service

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
)

type StoreOrderService interface {
	GetOrCreateOrder(sessionID uint, storeID uint) (*response.StoreOrderResponse, error)
	UpdateOrder(orderID uint, items []request.UpdateOrderItemRequest) error
	SubmitOrder(orderId uint) error
	GetOrderDetail(OrderID uint) (*response.StoreOrderDetailResponse, error)
}
