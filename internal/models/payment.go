package models

import (
	"github.com/google/uuid"
)

type Payment struct {
	BaseModel
	OrderID      uuid.UUID      `gorm:"type:char(36);index:idx_payment_order;not null" json:"order_id"`
	ChangeAmount float64        `json:"change_amount"`
	Status       string         `gorm:"type:varchar(20);index:idx_payment_status" json:"status"`
	Splits       []PaymentSplit `gorm:"foreignKey:PaymentID;constraint:OnDelete:CASCADE" json:"splits,omitempty"`
	SnapToken    string         `gorm:"-" json:"snap_token,omitempty"`
	SnapURL      string         `gorm:"-" json:"snap_url,omitempty"`
}

type PaymentSplit struct {
	BaseModel
	PaymentID       uuid.UUID     `gorm:"type:char(36);index:idx_split_payment;not null" json:"payment_id"`
	PaymentMethod   PaymentMethod `gorm:"type:varchar(20);index:idx_split_method" json:"payment_method"`
	Amount          float64       `gorm:"not null" json:"amount"`
	ReferenceNumber string        `json:"reference_number,omitempty"`
}
