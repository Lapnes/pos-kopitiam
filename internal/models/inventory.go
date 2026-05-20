package models

import (
	"github.com/google/uuid"
)

type RawMaterial struct {
	BaseModel
	Name          string  `gorm:"type:varchar(191);not null;index:idx_rawmaterial_name" json:"name"`
	Unit          string  `gorm:"not null" json:"unit"`
	CostPerUnit   float64 `json:"cost_per_unit"`
	CurrentStock  float64 `gorm:"default:0" json:"current_stock"`
	MinStockLevel float64 `json:"min_stock_level"`
}

type Recipe struct {
	BaseModel
	MenuID      uuid.UUID          `gorm:"type:char(36);uniqueIndex;not null" json:"menu_id"`
	Ingredients []RecipeIngredient `gorm:"foreignKey:RecipeID;references:ID;constraint:OnDelete:CASCADE" json:"ingredients,omitempty"`
}

type RecipeIngredient struct {
	BaseModel
	RecipeID      uuid.UUID `gorm:"type:char(36);index:idx_ingredient_recipe;not null" json:"recipe_id"`
	RawMaterialID uuid.UUID `gorm:"type:char(36);index:idx_ingredient_material;not null" json:"raw_material_id"`
	Quantity      float64   `gorm:"not null" json:"quantity"`
}
