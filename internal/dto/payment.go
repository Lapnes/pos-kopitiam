package dto

type PaymentSplitRequest struct {
	PaymentMethod   string  `json:"payment_method" binding:"required,oneof=cash qris debit_card credit_card transfer e_wallet"`
	Amount          float64 `json:"amount" binding:"required,min=0"`
	ReferenceNumber string  `json:"reference_number"`
}

type ProcessPaymentRequest struct {
	ChangeAmount float64               `json:"change_amount"`
	Splits       []PaymentSplitRequest `json:"splits" binding:"required,min=1"`
}