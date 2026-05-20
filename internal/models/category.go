package models

type Category struct {
	BaseModel
	Name      string `gorm:"type:varchar(191);not null;uniqueIndex" json:"name"`
	SortOrder int    `gorm:"default:0" json:"sort_order"`
	Menus     []Menu `gorm:"many2many:menu_categories;" json:"menus,omitempty"`
}
