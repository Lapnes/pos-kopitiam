package models

import (
	"time"

	"github.com/google/uuid"
)

type Shift struct {
	BaseModel
	BranchID          uuid.UUID  `gorm:"type:char(36);index:idx_shift_branch;not null" json:"branch_id"`
	OpenedBy          uuid.UUID  `gorm:"type:char(36);index:idx_shift_opener;not null" json:"opened_by"`
	OpeningTime       time.Time  `gorm:"index:idx_shift_opening" json:"opening_time"`
	OpeningCash       float64    `json:"opening_cash"`
	ClosedBy          *uuid.UUID `gorm:"type:char(36);index:idx_shift_closer" json:"closed_by,omitempty"`
	ClosingTime       *time.Time `json:"closing_time,omitempty"`
	ActualClosingCash float64    `json:"actual_closing_cash"`
	Status            string     `gorm:"type:varchar(20);index:idx_shift_status;not null" json:"status"`
}

type CashMovement struct {
	BaseModel
	ShiftID    uuid.UUID        `gorm:"type:char(36);index:idx_cashmovement_shift;not null" json:"shift_id"`
	Type       CashMovementType `gorm:"type:varchar(10);index:idx_cashmovement_type;not null" json:"type"`
	Amount     float64          `gorm:"not null" json:"amount"`
	Reason     string           `gorm:"type:text" json:"reason"`
	RecordedBy uuid.UUID        `gorm:"type:char(36);index:idx_cashmovement_user;not null" json:"recorded_by"`
	RecordedAt time.Time        `gorm:"index:idx_cashmovement_date" json:"recorded_at"`
}
