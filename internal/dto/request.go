package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginPINRequest struct {
	PIN string `json:"pin" binding:"required,len=6"`
}

type CreateOrderRequest struct {
	TableID        string             `json:"table_id"`
	CustomerName   string             `json:"customer_name"`
	CustomerPhone  string             `json:"customer_phone"`
	OrderType      string             `json:"order_type" binding:"required,oneof=dine_in takeaway delivery online"`
	Notes          string             `json:"notes"`
	DiscountAmount float64            `json:"discount_amount"`
	Items          []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type OrderItemRequest struct {
	MenuID   string  `json:"menu_id" binding:"required"`
	MenuName string  `json:"menu_name" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
	Price    float64 `json:"price" binding:"required,min=0"`
	Notes    string  `json:"notes"`
}
