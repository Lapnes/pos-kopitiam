package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	base.ID = uuid.New()
	return nil
}

type Branch struct {
	BaseModel
	Name    string `gorm:"type:varchar(100);not null" json:"name"`
	Address string `gorm:"type:text" json:"address"`
	Phone   string `gorm:"type:varchar(20)" json:"phone"`
	IsMain  bool   `gorm:"default:false" json:"is_main"`
	Areas   []Area `gorm:"foreignKey:BranchID" json:"areas,omitempty"`
}

type Area struct {
	BaseModel
	BranchID uuid.UUID `gorm:"type:char(36);not null;index" json:"branch_id"`
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
	Tables   []Table   `gorm:"foreignKey:AreaID" json:"tables,omitempty"`
}

type Table struct {
	BaseModel
	AreaID   uuid.UUID   `gorm:"type:char(36);not null;index" json:"area_id"`
	Name     string      `gorm:"type:varchar(50);not null" json:"name"`
	Capacity int         `gorm:"not null" json:"capacity"`
	Status   TableStatus `gorm:"type:varchar(20);default:'available'" json:"status"`
	QRCode   string      `gorm:"type:varchar(255)" json:"qr_code"`
}

// TASK 2.1: Removed Station field from Category (Station belongs to Menu now)
type Category struct {
	BaseModel
	Name      string `gorm:"type:varchar(100);not null" json:"name"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
}

// TASK 2.1 & 2.2: Removed CategoryID (M2M via menu_categories). Added IsRecipeBased.
type Menu struct {
	BaseModel
	Name          string     `gorm:"type:varchar(100);not null" json:"name"`
	Description   string     `gorm:"type:text" json:"description"`
	Price         float64    `gorm:"type:decimal(15,2);not null" json:"price"`
	CostPrice     float64    `gorm:"type:decimal(15,2);not null" json:"cost_price"` // HPP
	DailyStock    int        `gorm:"not null;default:0" json:"daily_stock"`
	IsActive      bool       `gorm:"default:true" json:"is_active"`
	IsRecipeBased bool       `gorm:"default:false" json:"is_recipe_based"` // TASK 2.2: Inventory branch flag
	Station       Station    `gorm:"type:varchar(20);not null" json:"station"`
	ImageURL      string     `gorm:"type:varchar(255)" json:"image_url"`
	Categories    []Category `gorm:"many2many:menu_categories;" json:"categories"` // TASK 2.1: M2M
}

type Employee struct {
	BaseModel
	BranchID uuid.UUID `gorm:"type:char(36);not null;index" json:"branch_id"`
	Name     string    `gorm:"type:varchar(100);not null" json:"name"`
	Email    string    `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password string    `gorm:"type:varchar(255);not null" json:"-"`
	PINCode  string    `gorm:"type:varchar(10);not null" json:"-"`
	Role     Role      `gorm:"type:varchar(20);not null" json:"role"`
	IsActive bool      `gorm:"default:true" json:"is_active"`
}

type Shift struct {
	BaseModel
	BranchID          uuid.UUID  `gorm:"type:char(36);not null;index" json:"branch_id"`
	OpenedBy          uuid.UUID  `gorm:"type:char(36);not null;index" json:"opened_by"`
	ClosedBy          *uuid.UUID `gorm:"type:char(36);index" json:"closed_by"`
	OpeningTime       time.Time  `gorm:"not null" json:"opening_time"`
	ClosingTime       *time.Time `json:"closing_time"`
	OpeningCash       float64    `gorm:"type:decimal(15,2);not null" json:"opening_cash"`
	ActualClosingCash float64    `gorm:"type:decimal(15,2)" json:"actual_closing_cash"`
	ExpectedCash      float64    `gorm:"type:decimal(15,2)" json:"expected_cash"`
	Variance          float64    `gorm:"type:decimal(15,2)" json:"variance"`
	TotalSales        float64    `gorm:"type:decimal(15,2)" json:"total_sales"`
	TotalTransactions int        `json:"total_transactions"`
	Status            string     `gorm:"type:varchar(20);default:'open'" json:"status"`
}

// TASK 2.5: Composite index on Status + CreatedAt
// TASK 1: Removed BigcapitalSyncID and SyncedToBigcapital
type Order struct {
	BaseModel
	OrderNumber   string      `gorm:"type:varchar(50);uniqueIndex;not null" json:"order_number"` // ORD-YYYYMMDD-XXXX
	BranchID      uuid.UUID   `gorm:"type:char(36);not null;index" json:"branch_id"`
	ShiftID       *uuid.UUID  `gorm:"type:char(36);index" json:"shift_id"`
	TableID       *uuid.UUID  `gorm:"type:char(36);index" json:"table_id"` // Nullable for takeaway/delivery
	EmployeeID    uuid.UUID   `gorm:"type:char(36);not null;index" json:"employee_id"`
	CustomerName  string      `gorm:"type:varchar(100)" json:"customer_name"`
	CustomerPhone string      `gorm:"type:varchar(20)" json:"customer_phone"`
	OrderType     OrderType   `gorm:"type:varchar(20);not null" json:"order_type"`
	Status        OrderStatus `gorm:"type:varchar(20);not null;default:'pending';index:idx_status_date" json:"status"`

	Subtotal       float64 `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	TaxAmount      float64 `gorm:"type:decimal(15,2);not null" json:"tax_amount"` // PPN
	ServiceCharge  float64 `gorm:"type:decimal(15,2);not null" json:"service_charge"`
	DiscountAmount float64 `gorm:"type:decimal(15,2);not null" json:"discount_amount"`
	Total          float64 `gorm:"type:decimal(15,2);not null" json:"total"`

	Notes string `gorm:"type:text" json:"notes"`

	// Cancellation tracking
	CancelledBy  *uuid.UUID `gorm:"type:char(36);index" json:"cancelled_by"`
	CancelledAt  *time.Time `json:"cancelled_at"`
	CancelReason string     `gorm:"type:text" json:"cancel_reason"`

	// Return tracking
	HasReturns bool `gorm:"default:false" json:"has_returns"`

	OrderDetails []OrderDetail `gorm:"foreignKey:OrderID" json:"order_details"`
	Payments     []Payment     `gorm:"foreignKey:OrderID" json:"payments"`
}

// TASK 3: OOP Fat Model — CalculateTotal encapsulates business rules on Order
func (o *Order) CalculateTotal() {
	var subtotal float64
	for _, item := range o.OrderDetails {
		if !item.IsVoided {
			subtotal += item.Subtotal
		}
	}
	o.Subtotal = subtotal
	o.TaxAmount = subtotal * 0.11     // 11% Tax
	o.ServiceCharge = subtotal * 0.05 // 5% Service Charge
	o.Total = (o.Subtotal + o.TaxAmount + o.ServiceCharge) - o.DiscountAmount
	if o.Total < 0 {
		o.Total = 0
	}
}

// TASK 2.3: Added CostPrice for historical COGS integrity
type OrderDetail struct {
	BaseModel
	OrderID   uuid.UUID `gorm:"type:char(36);not null;index" json:"order_id"`
	MenuID    uuid.UUID `gorm:"type:char(36);not null;index" json:"menu_id"`
	MenuName  string    `gorm:"type:varchar(100);not null" json:"menu_name"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	Price     float64   `gorm:"type:decimal(15,2);not null" json:"price"`
	CostPrice float64   `gorm:"type:decimal(15,2);not null" json:"cost_price"` // TASK 2.3: Snapshot HPP
	Subtotal  float64   `gorm:"type:decimal(15,2);not null" json:"subtotal"`
	Notes     string    `gorm:"type:text" json:"notes"`

	// Void tracking
	IsVoided   bool       `gorm:"default:false" json:"is_voided"`
	VoidedBy   *uuid.UUID `gorm:"type:char(36);index" json:"voided_by"`
	VoidedAt   *time.Time `json:"voided_at"`
	VoidReason string     `gorm:"type:text" json:"void_reason"`

	// KDS Tracking
	Station   Station   `gorm:"type:varchar(20);not null" json:"station"`
	KDSStatus KDSStatus `gorm:"type:varchar(20);default:'pending'" json:"kds_status"`
	Priority  int       `gorm:"default:0" json:"priority"` // 0: Normal, 1: Urgent, 2: VIP
}

// TASK 2.4: Removed PaymentMethod, Amount, CashTendered from Payment header.
// Payment now represents the full payment event; splits contain per-method detail.
type Payment struct {
	BaseModel
	OrderID         uuid.UUID  `gorm:"type:char(36);not null;index" json:"order_id"`
	ShiftID         *uuid.UUID `gorm:"type:char(36);index" json:"shift_id"`
	TotalPaid       float64    `gorm:"type:decimal(15,2);not null" json:"total_paid"` // TASK 2.4
	ChangeAmount    float64    `gorm:"type:decimal(15,2)" json:"change_amount"`
	ReferenceNumber string     `gorm:"type:varchar(100)" json:"reference_number"`
	Status          string     `gorm:"type:varchar(20);default:'success'" json:"status"` // success, refunded

	// Refund tracking
	RefundedBy   *uuid.UUID `gorm:"type:char(36);index" json:"refunded_by"`
	RefundReason string     `gorm:"type:text" json:"refund_reason"`

	Splits []PaymentSplit `gorm:"foreignKey:PaymentID" json:"splits,omitempty"`
}

type PaymentSplit struct {
	BaseModel
	PaymentID       uuid.UUID     `gorm:"type:char(36);not null;index" json:"payment_id"`
	PaymentMethod   PaymentMethod `gorm:"type:varchar(20);not null" json:"payment_method"`
	Amount          float64       `gorm:"type:decimal(15,2);not null" json:"amount"`
	ReferenceNumber string        `gorm:"type:varchar(100)" json:"reference_number"`
}

type Discount struct {
	BaseModel
	Name           string       `gorm:"type:varchar(100);not null" json:"name"`
	Type           DiscountType `gorm:"type:varchar(20);not null" json:"type"`
	Value          float64      `gorm:"type:decimal(15,2);not null" json:"value"`
	MinOrderAmount float64      `gorm:"type:decimal(15,2)" json:"min_order_amount"`
	MaxDiscount    float64      `gorm:"type:decimal(15,2)" json:"max_discount"`
	StartDate      *time.Time   `json:"start_date"`
	EndDate        *time.Time   `json:"end_date"`
	IsActive       bool         `gorm:"default:true" json:"is_active"`
}

type DiscountApplication struct {
	BaseModel
	OrderID        uuid.UUID `gorm:"type:char(36);not null;index" json:"order_id"`
	DiscountID     uuid.UUID `gorm:"type:char(36);not null;index" json:"discount_id"`
	DiscountAmount float64   `gorm:"type:decimal(15,2);not null" json:"discount_amount"`
}

type Reservation struct {
	BaseModel
	BranchID        uuid.UUID `gorm:"type:char(36);not null;index" json:"branch_id"`
	TableID         uuid.UUID `gorm:"type:char(36);not null;index" json:"table_id"`
	CustomerName    string    `gorm:"type:varchar(100);not null" json:"customer_name"`
	CustomerPhone   string    `gorm:"type:varchar(20);not null" json:"customer_phone"`
	ReservationDate time.Time `gorm:"not null" json:"reservation_date"`
	GuestCount      int       `gorm:"not null" json:"guest_count"`
	DepositAmount   float64   `gorm:"type:decimal(15,2);default:0" json:"deposit_amount"`
	Status          string    `gorm:"type:varchar(20);default:'confirmed'" json:"status"` // confirmed, cancelled, completed, no_show
}

type AuditLog struct {
	BaseModel
	TableName   string     `gorm:"type:varchar(50);not null;index" json:"table_name"`
	RecordID    string     `gorm:"type:varchar(36);not null;index" json:"record_id"`
	Action      string     `gorm:"type:varchar(20);not null" json:"action"` // CREATE, UPDATE, DELETE, VOID
	OldValues   string     `gorm:"type:json" json:"old_values"`
	NewValues   string     `gorm:"type:json" json:"new_values"`
	PerformedBy *uuid.UUID `gorm:"type:char(36);index" json:"performed_by"`
	IPAddress   string     `gorm:"type:varchar(50)" json:"ip_address"`
	UserAgent   string     `gorm:"type:text" json:"user_agent"`
}

// --- ADVANCED INVENTORY & BOM ---

type RawMaterial struct {
	BaseModel
	Name         string  `gorm:"type:varchar(100);not null" json:"name"`
	Unit         string  `gorm:"type:varchar(20);not null" json:"unit"` // e.g., gram, ml, pcs
	CurrentStock float64 `gorm:"type:decimal(15,2);not null;default:0" json:"current_stock"`
	MinimumStock float64 `gorm:"type:decimal(15,2);not null;default:0" json:"minimum_stock"`
	CostPerUnit  float64 `gorm:"type:decimal(15,2);not null" json:"cost_per_unit"`
}

// FIX: Changed uniqueIndex to composite unique index that includes deleted_at
// This allows soft-deleted recipes to be recreated for the same menu
type Recipe struct {
	BaseModel
	MenuID       uuid.UUID    `gorm:"type:char(36);not null;uniqueIndex:idx_recipe_menu_deleted" json:"menu_id"`
	Instructions string       `gorm:"type:text" json:"instructions"`
	Ingredients  []RecipeItem `gorm:"foreignKey:RecipeID" json:"ingredients"`
}

type RecipeItem struct {
	BaseModel
	RecipeID      uuid.UUID `gorm:"type:char(36);not null;index" json:"recipe_id"`
	RawMaterialID uuid.UUID `gorm:"type:char(36);not null;index" json:"raw_material_id"`
	Quantity      float64   `gorm:"type:decimal(15,2);not null" json:"quantity"` // Amount needed per 1 portion
}

// --- RETURN & REFUND SYSTEM ---

type OrderReturn struct {
	BaseModel
	OrderID        uuid.UUID  `gorm:"type:char(36);not null;index" json:"order_id"`
	OrderDetailID  *uuid.UUID `gorm:"type:char(36);index" json:"order_detail_id"`       // Nullable if returning whole order
	ProcessedBy    uuid.UUID  `gorm:"type:char(36);not null;index" json:"processed_by"` // Manager ID
	Reason         string     `gorm:"type:text;not null" json:"reason"`
	OriginalAmount float64    `gorm:"type:decimal(15,2);not null" json:"original_amount"`
	ReturnAmount   float64    `gorm:"type:decimal(15,2);not null" json:"return_amount"` // Max 80%
}
