package models

import (
	"time"

	"github.com/google/uuid"
)

type Discount struct {
	BaseModel
	Name           string       `gorm:"type:varchar(191);not null;index:idx_discount_name" json:"name"`
	Type           DiscountType `gorm:"type:varchar(20);index:idx_discount_type;not null" json:"type"`
	Value          float64      `gorm:"not null" json:"value"`
	MinOrderAmount float64      `json:"min_order_amount"`
	MaxDiscount    float64      `json:"max_discount"`
	StartDate      *time.Time   `json:"start_date,omitempty"`
	EndDate        *time.Time   `json:"end_date,omitempty"`
	IsActive       bool         `gorm:"default:true;index:idx_discount_active" json:"is_active"`
}

type DiscountApplication struct {
	BaseModel
	OrderID        uuid.UUID `gorm:"type:char(36);index:idx_discountapp_order;not null" json:"order_id"`
	DiscountID     uuid.UUID `gorm:"type:char(36);index:idx_discountapp_discount;not null" json:"discount_id"`
	DiscountAmount float64   `json:"discount_amount"`
}
