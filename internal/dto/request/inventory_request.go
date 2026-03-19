package request

type UpdateInventoryRequest struct {
	ProductID uint    `json:"productId"`
	Quantity  float64 `json:"quantity"`
}

type BatchAdjustInventoryRequest struct {
	Adjustments []AdjustmentItem `json:"adjustments" binding:"required,min=1,dive"`
}

type AdjustmentItem struct {
	StoreID   uint    `json:"store_id" binding:"required"`
	ProductID uint    `json:"product_id" binding:"required"`
	Delta     float64 `json:"delta" binding:"required"`
	Reason    string  `json:"reason"`
}
