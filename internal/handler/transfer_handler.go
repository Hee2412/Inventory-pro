package handler

import (
	"Inventory-pro/internal/dto/request"
	"Inventory-pro/internal/dto/response"
	"Inventory-pro/internal/service"
	"Inventory-pro/pkg/pagination"
	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	service service.TransferService
}

func NewTransferHandler(service service.TransferService) *TransferHandler {
	return &TransferHandler{service: service}
}

// CreateTransfer POST /api/store/transfers
func (h *TransferHandler) CreateTransfer(c *gin.Context) {
	var req request.CreateTransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	storeID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := h.service.CreateTransfer(req, storeID, storeID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, result)
}

// AdminCreateTransfer POST /api/admin/transfers
func (h *TransferHandler) AdminCreateTransfer(c *gin.Context) {
	var req struct {
		FromStoreID uint `json:"from_store_id" binding:"required"`
		request.CreateTransferRequest
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	adminID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := h.service.CreateTransfer(req.CreateTransferRequest, adminID, req.FromStoreID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Created(c, result)
}

// GetTransferDetail GET /api/store/transfers/:transferId
//
//	GET /api/admin/transfers/:transferId
func (h *TransferHandler) GetTransferDetail(c *gin.Context) {
	transferID, err := getIDParam(c, "transferId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := h.service.GetTransferDetail(transferID)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetMyTransfers GET /api/store/transfers?direction=in|out
func (h *TransferHandler) GetMyTransfers(c *gin.Context) {
	storeID, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	direction := c.Query("direction") // "in" | "out" | "" (all)
	result, err := h.service.GetMyTransfers(storeID, direction)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// GetAllTransfers GET /api/admin/transfers?page=1&limit=20&status=PENDING
func (h *TransferHandler) GetAllTransfers(c *gin.Context) {
	var params request.TransferSearchParams
	if err := c.ShouldBindQuery(&params); err != nil {
		response.HandleError(c, err)
		return
	}
	params.SetDefaults()
	orders, total, err := h.service.GetAllTransfers(params.Page, params.Limit, params.Status)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	paginatedResponse := pagination.NewResponse(orders, params.Page, params.Limit, total)
	response.Success(c, paginatedResponse)
}

// ApproveTransfer POST /api/store/transfers/:transferId/approve  (ToStore xác nhận)
//
//	POST /api/admin/transfers/:transferId/approve  (Admin xác nhận)
func (h *TransferHandler) ApproveTransfer(c *gin.Context) {
	transferID, err := getIDParam(c, "transferId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	approvedBy, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	result, err := h.service.ApproveTransfer(transferID, approvedBy)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	response.Success(c, result)
}

// CancelTransfer POST /api/store/transfers/:transferId/cancel
//
//	POST /api/admin/transfers/:transferId/cancel
func (h *TransferHandler) CancelTransfer(c *gin.Context) {
	transferID, err := getIDParam(c, "transferId")
	if err != nil {
		response.HandleError(c, err)
		return
	}
	cancelledBy, err := getUserID(c)
	if err != nil {
		response.HandleError(c, err)
		return
	}
	var req request.CancelTransferRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		response.HandleError(c, err)
		return
	}
	if err = h.service.CancelTransfer(transferID, cancelledBy, req.Reason); err != nil {
		response.HandleError(c, err)
		return
	}
	response.Message(c, "Transfer cancelled successfully")
}
