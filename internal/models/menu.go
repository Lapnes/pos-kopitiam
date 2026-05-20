package models

type Menu struct {
	BaseModel
	Name          string     `gorm:"type:varchar(191);not null;index:idx_menu_name" json:"name"`
	Price         float64    `gorm:"not null" json:"price"`
	IsActive      bool       `gorm:"default:true;index:idx_menu_active" json:"is_active"`
	IsRecipeBased bool       `gorm:"default:false" json:"is_recipe_based"`
	DailyStock    int        `gorm:"default:0" json:"daily_stock"`
	Station       string     `gorm:"type:varchar(50);index:idx_menu_station" json:"station"`
	Categories    []Category `gorm:"many2many:menu_categories;" json:"categories,omitempty"`
}
