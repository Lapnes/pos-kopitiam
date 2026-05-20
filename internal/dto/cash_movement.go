package dto

type CashMovementRequest struct {
	Type   string  `json:"type" binding:"required,oneof=in out"`
	Amount float64 `json:"amount" binding:"required,min=0"`
	Reason string  `json:"reason" binding:"required"`
}