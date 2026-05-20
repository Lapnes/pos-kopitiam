package models

import (
	"github.com/google/uuid"
)

type StockAdjustment struct {
	BaseModel
	RawMaterialID  uuid.UUID `gorm:"type:char(36);index:idx_adjustment_material;not null" json:"raw_material_id"`
	QuantityBefore float64   `json:"quantity_before"`
	QuantityAfter  float64   `gorm:"not null" json:"quantity_after"`
	Reason         string    `gorm:"type:text;not null" json:"reason"`
	AdjustedBy     uuid.UUID `gorm:"type:char(36);index:idx_adjustment_user;not null" json:"adjusted_by"`
}
