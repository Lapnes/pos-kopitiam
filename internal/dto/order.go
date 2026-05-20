package dto

type OrderItemRequest struct {
	MenuID   string  `json:"menu_id" binding:"required"`
	MenuName string  `json:"menu_name" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
	Price    float64 `json:"price" binding:"required,min=0"`
	Notes    string  `json:"notes"`
}

type CreateOrderRequest struct {
	TableID        string             `json:"table_id"`
	CustomerID     string             `json:"customer_id"`
	OrderType      string             `json:"order_type" binding:"required,oneof=dine_in takeaway delivery online"`
	Notes          string             `json:"notes"`
	DiscountAmount float64            `json:"discount_amount"`
	Items          []OrderItemRequest `json:"items" binding:"required,min=1"`
}

type UpdateOrderNotesRequest struct {
	Notes string `json:"notes"`
}

type VoidItemRequest struct {
	OrderDetailID string `json:"order_detail_id" binding:"required"`
	Reason        string `json:"reason" binding:"required"`
}

type VoidOrderRequest struct {
	Reason string `json:"reason" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}