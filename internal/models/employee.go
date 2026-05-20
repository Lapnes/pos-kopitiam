package models

import (
	"github.com/google/uuid"
)

type User struct {
	BaseModel
	Name         string    `gorm:"type:varchar(191);not null;index:idx_user_name" json:"name"`
	Email        string    `gorm:"type:varchar(191);uniqueIndex;not null" json:"email"`
	Phone        string    `gorm:"type:varchar(191);index:idx_user_phone" json:"phone"`
	PasswordHash string    `gorm:"not null" json:"-"`
	PIN          string    `gorm:"not null" json:"-"`
	Role         Role      `gorm:"type:varchar(20);index:idx_user_role;not null" json:"role"`
	BranchID     uuid.UUID `gorm:"type:char(36);index:idx_user_branch" json:"branch_id"`
	IsActive     bool      `gorm:"default:true;index:idx_user_active" json:"is_active"`
}
