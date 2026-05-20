package dto

type OpenShiftRequest struct {
	OpeningCash float64 `json:"opening_cash" binding:"required,min=0"`
}

type CloseShiftRequest struct {
	ActualClosingCash float64 `json:"actual_closing_cash" binding:"required,min=0"`
}