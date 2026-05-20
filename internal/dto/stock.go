package dto

type StockAdjustmentRequest struct {
	RawMaterialID string  `json:"raw_material_id" binding:"required"`
	QuantityAfter float64 `json:"quantity_after" binding:"required"`
	Reason        string  `json:"reason" binding:"required"`
}