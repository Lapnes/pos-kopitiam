package models

import (
	"github.com/google/uuid"
)

type Order struct {
	BaseModel
	OrderNumber    string        `gorm:"type:varchar(191);uniqueIndex;not null" json:"order_number"`
	BranchID       uuid.UUID     `gorm:"type:char(36);index:idx_orders_branch_status;index:idx_orders_branch_created" json:"branch_id"`
	TableID        *uuid.UUID    `gorm:"type:char(36);index:idx_orders_table" json:"table_id,omitempty"`
	EmployeeID     uuid.UUID     `gorm:"type:char(36);index:idx_orders_employee" json:"employee_id"`
	CustomerID     *uuid.UUID    `gorm:"type:char(36);index:idx_orders_customer" json:"customer_id,omitempty"`
	ShiftID        *uuid.UUID    `gorm:"type:char(36);index:idx_orders_shift" json:"shift_id,omitempty"`
	OrderType      OrderType     `gorm:"type:varchar(20);index:idx_orders_type;not null" json:"order_type"`
	Status         OrderStatus   `gorm:"type:varchar(20);index:idx_orders_branch_status;index:idx_orders_status;not null" json:"status"`
	Subtotal       float64       `json:"subtotal"` // 3NF — stored for historical accuracy
	TaxAmount      float64       `json:"tax_amount"`
	ServiceCharge  float64       `json:"service_charge"`
	DiscountAmount float64       `json:"discount_amount"`
	Total          float64       `json:"total"`
	Notes          string        `gorm:"type:text" json:"notes"`
	HasReturns     bool          `gorm:"default:false" json:"has_returns"`
	OrderDetails   []OrderDetail `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"order_details,omitempty"`
	Payments       []Payment     `gorm:"foreignKey:OrderID" json:"payments,omitempty"`
}

type OrderDetail struct {
	BaseModel
	OrderID    uuid.UUID `gorm:"type:char(36);index:idx_details_order;not null" json:"order_id"`
	MenuID     uuid.UUID `gorm:"type:char(36);index:idx_details_menu;not null" json:"menu_id"`
	MenuName   string    `gorm:"not null" json:"menu_name"`
	Quantity   int       `gorm:"not null" json:"quantity"`
	Price      float64   `gorm:"not null" json:"price"`
	CostPrice  float64   `json:"cost_price"`
	Subtotal   float64   `gorm:"not null" json:"subtotal"`
	Notes      string    `gorm:"type:text" json:"notes"`
	Station    string    `gorm:"type:varchar(50)" json:"station"`
	IsVoided   bool      `gorm:"default:false;index:idx_details_voided" json:"is_voided"`
	VoidReason string    `json:"void_reason,omitempty"`
}
