package models

import (
	"time"

	"github.com/google/uuid"
)

type TaxConfig struct {
	BaseModel
	BranchID      *uuid.UUID `gorm:"type:char(36);index:idx_tax_branch" json:"branch_id,omitempty"`
	Name          string     `gorm:"type:varchar(191);not null;index:idx_tax_name" json:"name"`
	TaxType       TaxType    `gorm:"type:varchar(20);index:idx_tax_type;not null" json:"tax_type"`
	Percentage    float64    `gorm:"not null" json:"percentage"`
	IsActive      bool       `gorm:"default:true;index:idx_tax_active" json:"is_active"`
	EffectiveFrom time.Time  `gorm:"index:idx_tax_effective" json:"effective_from"`
}
