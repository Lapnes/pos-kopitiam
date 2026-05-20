package dto

type ProcessReturnRequest struct {
	ReturnAmount float64 `json:"return_amount" binding:"required,min=0"`
	Reason       string  `json:"reason" binding:"required"`
}