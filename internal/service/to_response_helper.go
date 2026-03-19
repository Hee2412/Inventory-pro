package service

import (
	"Inventory-pro/internal/domain"
	"Inventory-pro/internal/dto/response"
	"time"
)

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

func toProductResponse(product *domain.Product) *response.ProductResponse {
	return &response.ProductResponse{
		ID:          product.ID,
		ProductName: product.ProductName,
		ProductCode: product.ProductCode,
		Unit:        product.Unit,
		MOQ:         product.MOQ,
		OM:          product.OM,
		Type:        product.Type,
		OrderCycle:  product.OrderCycle,
		AuditCycle:  product.AuditCycle,
		IsActive:    product.IsActive,
	}
}

func toUserResponse(user *domain.User) *response.UserResponse {
	var lastLogin string
	if user.LastLogin != nil {
		lastLogin = user.LastLogin.Format(time.DateTime)
	} else {
		lastLogin = ""
	}
	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Role:      user.Role,
		StoreName: user.StoreName,
		StoreCode: user.StoreCode,
		IsActive:  user.IsActive,
		CreateAt:  user.CreatedAt,
		LastLogin: lastLogin,
		DeletedAt: user.DeletedAt,
	}
}

func toAuditSessionResponse(session *domain.AuditSession) *response.AuditSessionResponse {
	return &response.AuditSessionResponse{
		SessionID: session.ID,
		Title:     session.Title,
		AuditType: session.AuditType,
		StartDate: session.StartDate,
		EndDate:   session.EndDate,
		Status:    session.Status,
		CreatedBy: session.CreatedBy,
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

func toOrderSessionResponse(order *domain.OrderSession) *response.OrderSessionResponse {
	return &response.OrderSessionResponse{
		ID:         order.ID,
		Title:      order.Title,
		OrderCycle: order.OrderCycle,
		Status:     order.Status,
		Deadline:   order.Deadline,
		DeliveryAt: order.DeliveryDate,
		CreatedAt:  order.CreatedAt,
	}
}

func toAuditReportItemResponse(items []*domain.StoreAuditReport) []*response.AuditItemsResponse {
	result := make([]*response.AuditItemsResponse, 0, len(items))
	for _, item := range items {
		result = append(result, &response.AuditItemsResponse{
			ProductID:   item.ProductID,
			ProductName: item.Product.ProductName,
			SystemStock: item.SystemStock,
			ActualStock: item.ActualStock,
			Variance:    item.Variance,
		})
	}
	return result
}

func toInventoryResponse(inv *domain.StoreInventory) *response.InventoryResponse {
	resp := &response.InventoryResponse{
		ID:        inv.ID,
		StoreID:   inv.StoreID,
		ProductID: inv.ProductID,
		Quantity:  inv.Quantity,
		UpdatedBy: inv.UpdatedBy,
		UpdatedAt: inv.UpdatedAt,
	}
	if inv.Store != nil {
		resp.StoreName = inv.Store.StoreName
	}
	if inv.Product != nil {
		resp.ProductName = inv.Product.ProductName
		resp.ProductCode = inv.Product.ProductCode
	}
	return resp
}

func toTransferOrderResponse(order *domain.TransferOrder) *response.TransferOrderResponse {
	resp := &response.TransferOrderResponse{
		ID:           order.ID,
		FromStoreID:  order.FromStoreID,
		ToStoreID:    order.ToStoreID,
		Status:       order.Status,
		Note:         order.Note,
		CreatedBy:    order.CreatedBy,
		CreatedAt:    order.CreatedAt,
		ApprovedAt:   order.ApprovedAt,
		CancelledAt:  order.CancelledAt,
		CancelReason: order.CancelReason,
	}
	return resp
}

func toTransferItemResponse(items []*domain.TransferOrderItem) []response.TransferItemResponse {
	result := make([]response.TransferItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, response.TransferItemResponse{
			ID:          item.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}
	return result
}
