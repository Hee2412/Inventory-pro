package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"github.com/gin-gonic/gin"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: inventoryService}
}

// GetStoreInventory /api/admin/inventory/stores/:storeId
func (h *InventoryHandler) GetStoreInventory(c *gin.Context) {
	storeID, err := getIDParam(c, "storeId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := h.service.GetStoreInventory(storeID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// UpdateInventory /api/inventory/admin/stores/:storeId
func (h *InventoryHandler) UpdateInventory(c *gin.Context) {
	storeID, err := getIDParam(c, "storeId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	var req request.UpdateInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	adminID, err := getIDParam(c, "adminId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	err = h.service.UpdateInventory(storeID, req, adminID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Inventory updated")
}

// AdjustInventory POST /api/admin/inventory/adjust
func (h *InventoryHandler) AdjustInventory(c *gin.Context) {
	adminID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	var req request.BatchAdjustInventoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}

	err = h.service.BatchAdjustInventory(req.Adjustments, adminID)
	if err != nil {
		response.HandleError(c, err)
		return
	}

	response.Message(c, "Inventory adjusted successfully")
}

// GetAllInventory /api/admin/inventory
func (h *InventoryHandler) GetAllInventory(c *gin.Context) {
	result, err := h.service.GetAllInventory()
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetProductInventoryAcrossStore /api/admin/inventory/products/:productId
func (h *InventoryHandler) GetProductInventoryAcrossStore(c *gin.Context) {
	productID, err := getIDParam(c, "productId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := h.service.GetProductInventoryAcrossStores(productID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}
