package models

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	BaseModel
	UserID    uuid.UUID `gorm:"type:char(36);index:idx_audit_user;not null" json:"user_id"`
	Action    string    `gorm:"type:varchar(50);index:idx_audit_action;not null" json:"action"`
	Entity    string    `gorm:"type:varchar(50);index:idx_audit_entity;not null" json:"entity"`
	EntityID  uuid.UUID `gorm:"type:char(36);index:idx_audit_entity_id" json:"entity_id"`
	OldValue  string    `gorm:"type:json" json:"old_value,omitempty"`
	NewValue  string    `gorm:"type:json" json:"new_value,omitempty"`
	Timestamp time.Time `gorm:"index:idx_audit_timestamp;not null" json:"timestamp"`
}
