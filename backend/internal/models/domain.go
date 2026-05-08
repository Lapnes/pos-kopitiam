package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ============================================
// MASTER DATA
// ============================================

// Category (1) -----> (*) Menu
type Category struct {
	CategoryID   uint64 `gorm:"primaryKey;autoIncrement" json:"category_id"`
	CategoryName string `gorm:"type:varchar(100);not null" json:"category_name"`

	// Relasi tanpa explicit FK constraint
	Menus []Menu `json:"menus,omitempty"`
}

// Menu (*) -----> (1) Category
// Menu (1) -----> (*) OrderDetail
type Menu struct {
	MenuID     uint64 `gorm:"primaryKey;autoIncrement" json:"menu_id"`
	CategoryID uint64 `gorm:"not null" json:"category_id"`
	MenuName   string `gorm:"type:varchar(150);not null" json:"menu_name"`
	Price      int64  `gorm:"not null" json:"price"`
	DailyStock int64  `gorm:"default:0" json:"daily_stock"`

	// Relasi tanpa explicit FK constraint
	Category     Category      `json:"category,omitempty"`
	OrderDetails []OrderDetail `json:"order_details,omitempty"`
}

// Employee (1) -----> (*) Order
type Employee struct {
	EmployeeID   uint64 `gorm:"primaryKey;autoIncrement; not null" json:"employee_id"`
	EmployeeName string `gorm:"type:varchar(150);not null" json:"employee_name"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	PhoneNumber  string `gorm:"type:varchar(20)" json:"phone_number"`

	// Relasi tanpa explicit FK constraint
	Orders []Order `json:"orders,omitempty"`
}

// ============================================
// TRANSACTION DATA
// ============================================

// Order (*) -----> (0..1) Employee
// Order (1) -----> (*) OrderDetail
type Order struct {
	OrderID    string  `gorm:"type:char(36);primaryKey" json:"order_id"`
	EmployeeID *uint64 `json:"employee_id"` // Nullable
	TotalPrice int64   `gorm:"not null;default:0" json:"total_price"`

	// Relasi tanpa explicit FK constraint
	Employee     *Employee     `json:"employee,omitempty"`
	OrderDetails []OrderDetail `json:"order_details,omitempty"`
}

func (o *Order) BeforeCreate(tx *gorm.DB) (err error) {
	if o.OrderID == "" {
		o.OrderID = uuid.New().String()
	}
	return
}

// OrderDetail (*) -----> (1) Order
// OrderDetail (*) -----> (1) Menu
type OrderDetail struct {
	OrderDetailID string `gorm:"type:char(36);primaryKey" json:"order_detail_id"`
	OrderID       string `gorm:"type:char(36);not null" json:"order_id"`
	MenuID        uint64 `gorm:"not null" json:"menu_id"`
	Quantity      int    `gorm:"not null" json:"quantity"`
	UnitPrice     int64  `gorm:"not null" json:"unit_price"`
	Subtotal      int64  `gorm:"not null" json:"subtotal"`

	// Relasi tanpa explicit FK constraint
	Order Order `json:"order,omitempty"`
	Menu  Menu  `json:"menu,omitempty"`
}

func (od *OrderDetail) BeforeCreate(tx *gorm.DB) (err error) {
	if od.OrderDetailID == "" {
		od.OrderDetailID = uuid.New().String()
	}
	od.Subtotal = int64(od.Quantity) * od.UnitPrice
	return
}
