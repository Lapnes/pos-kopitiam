package models

import (
	"time"

	"github.com/google/uuid"
)

type OrderReturn struct {
	BaseModel
	OrderID        uuid.UUID `gorm:"type:char(36);index:idx_return_order;not null" json:"order_id"`
	ProcessedBy    uuid.UUID `gorm:"type:char(36);index:idx_return_processor;not null" json:"processed_by"`
	Reason         string    `gorm:"type:text;not null" json:"reason"`
	OriginalAmount float64   `json:"original_amount"`
	ReturnAmount   float64   `gorm:"not null" json:"return_amount"`
	ProcessedAt    time.Time `gorm:"index:idx_return_date" json:"processed_at"`
}
