package dto

type OrderReportDTO struct {
	OrderID      string `json:"order_id"`
	EmployeeName string `json:"employee_name"`
	MenuName     string `json:"menu_name"`
	Quantity     int    `json:"quantity"`
	UnitPrice    int64  `json:"unit_price"`
	Subtotal     int64  `json:"subtotal"`
	TotalPrice   int64  `json:"total_price"`
}

type OrderResponseDTO struct {
	OrderID    string              `json:"order_id"`
	EmployeeID *uint64             `json:"employee_id"`
	TotalPrice int64               `json:"total_price"`
	Items      []OrderItemResponse `json:"items"`
}

type OrderItemResponse struct {
	MenuID    uint64 `json:"menu_id"`
	MenuName  string `json:"menu_name"`
	Quantity  int    `json:"quantity"`
	UnitPrice int64  `json:"unit_price"`
	Subtotal  int64  `json:"subtotal"`
}
