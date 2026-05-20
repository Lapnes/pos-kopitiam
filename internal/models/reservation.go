package models

import (
	"time"

	"github.com/google/uuid"
)

type Reservation struct {
	BaseModel
	BranchID        uuid.UUID         `gorm:"type:char(36);index:idx_reservation_branch;not null" json:"branch_id"`
	TableID         uuid.UUID         `gorm:"type:char(36);index:idx_reservation_table;not null" json:"table_id"`
	CustomerID      *uuid.UUID        `gorm:"type:char(36);index:idx_reservation_customer" json:"customer_id,omitempty"`
	CustomerName    string            `gorm:"not null" json:"customer_name"`
	CustomerPhone   string            `gorm:"not null" json:"customer_phone"`
	ReservationDate time.Time         `gorm:"index:idx_reservation_date" json:"reservation_date"`
	GuestCount      int               `gorm:"not null" json:"guest_count"`
	DepositAmount   float64           `json:"deposit_amount"`
	Status          ReservationStatus `gorm:"type:varchar(20);index:idx_reservation_status;not null" json:"status"`
}
