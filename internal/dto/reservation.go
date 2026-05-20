package dto

type ReservationRequest struct {
	TableID         string  `json:"table_id" binding:"required"`
	CustomerID      string  `json:"customer_id"`
	CustomerName    string  `json:"customer_name" binding:"required"`
	CustomerPhone   string  `json:"customer_phone" binding:"required"`
	ReservationDate string  `json:"reservation_date" binding:"required"` // YYYY-MM-DD HH:MM
	GuestCount      int     `json:"guest_count" binding:"required,min=1"`
	DepositAmount   float64 `json:"deposit_amount"`
}