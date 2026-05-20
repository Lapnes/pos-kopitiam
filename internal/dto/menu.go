package dto

type CreateMenuRequest struct {
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"required,min=0"`
	DailyStock    int      `json:"daily_stock" binding:"min=0"`
	IsActive      bool     `json:"is_active"`
	IsRecipeBased bool     `json:"is_recipe_based"`
	Station       string   `json:"station" binding:"required,oneof=kitchen bar pastry none"`
	ImageURL      string   `json:"image_url"`
	CategoryIDs   []string `json:"category_ids"`
}

type UpdateMenuRequest struct {
	Name          string   `json:"name"`
	Description   string   `json:"description"`
	Price         float64  `json:"price" binding:"omitempty,min=0"`
	DailyStock    *int     `json:"daily_stock" binding:"omitempty,min=0"`
	IsActive      *bool    `json:"is_active"`
	IsRecipeBased *bool    `json:"is_recipe_based"`
	Station       string   `json:"station" binding:"omitempty,oneof=kitchen bar pastry none"`
	ImageURL      string   `json:"image_url"`
	CategoryIDs   []string `json:"category_ids"`
}