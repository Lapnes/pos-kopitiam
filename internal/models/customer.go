package models

type Customer struct {
	BaseModel
	Name          string `gorm:"type:varchar(191);not null;index:idx_customer_name" json:"name"`
	Phone         string `gorm:"type:varchar(191);uniqueIndex;not null" json:"phone"`
	Email         string `gorm:"type:varchar(191);uniqueIndex" json:"email,omitempty"`
	LoyaltyPoints int    `gorm:"default:0" json:"loyalty_points"`
}
