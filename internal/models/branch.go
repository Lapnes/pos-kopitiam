package models

type Branch struct {
	BaseModel
	Name    string `gorm:"type:varchar(191);not null;index:idx_branch_name" json:"name"`
	Address string `json:"address"`
	Phone   string `gorm:"type:varchar(191);index:idx_branch_phone" json:"phone"`
}
