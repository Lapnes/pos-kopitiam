package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetDB() *gorm.DB

	// RawMaterial
	CreateRawMaterial(rm *models.RawMaterial) error
	UpdateRawMaterial(rm *models.RawMaterial) error
	GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error)
	GetRawMaterialsByIDs(ids []uuid.UUID) ([]models.RawMaterial, error) // BATCH — avoids N+1
	GetAllRawMaterials() ([]models.RawMaterial, error)
	DeleteRawMaterial(id uuid.UUID) error

	// Recipe
	CreateRecipe(recipe *models.Recipe) error
	UpdateRecipe(recipe *models.Recipe) error
	GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error)
	GetRecipeByMenuIDUnscoped(menuID uuid.UUID) (*models.Recipe, error)
	GetRecipeByID(id uuid.UUID) (*models.Recipe, error)
	GetRecipesByMenuIDs(menuIDs []uuid.UUID) ([]models.Recipe, error) // BATCH — avoids N+1
	DeleteRecipe(id uuid.UUID) error
	HardDeleteRecipe(id uuid.UUID) error
}

type inventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) InventoryRepository {
	return &inventoryRepository{db: db}
}

func (r *inventoryRepository) GetDB() *gorm.DB {
	return r.db
}

// RawMaterial
func (r *inventoryRepository) CreateRawMaterial(rm *models.RawMaterial) error {
	return r.db.Create(rm).Error
}

func (r *inventoryRepository) UpdateRawMaterial(rm *models.RawMaterial) error {
	return r.db.Omit("created_at").Save(rm).Error
}

func (r *inventoryRepository) GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error) {
	var rm models.RawMaterial
	err := r.db.First(&rm, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

// GetRawMaterialsByIDs — BATCH lookup: 1 query with IN clause instead of N queries
func (r *inventoryRepository) GetRawMaterialsByIDs(ids []uuid.UUID) ([]models.RawMaterial, error) {
	if len(ids) == 0 {
		return []models.RawMaterial{}, nil
	}
	var materials []models.RawMaterial
	err := r.db.Where("id IN ?", ids).Find(&materials).Error
	return materials, err
}

func (r *inventoryRepository) GetAllRawMaterials() ([]models.RawMaterial, error) {
	var materials []models.RawMaterial
	err := r.db.Find(&materials).Error
	return materials, err
}

func (r *inventoryRepository) DeleteRawMaterial(id uuid.UUID) error {
	return r.db.Delete(&models.RawMaterial{}, "id = ?", id).Error
}

// Recipe
// FIX: CreateRecipe now explicitly saves ingredients after creating the recipe
func (r *inventoryRepository) CreateRecipe(recipe *models.Recipe) error {
	// First create the recipe itself (omitting Ingredients to avoid GORM double-saving association)
	if err := r.db.Omit("Ingredients").Create(recipe).Error; err != nil {
		return err
	}

	// Then create ingredients with the recipe_id set
	if len(recipe.Ingredients) > 0 {
		for i := range recipe.Ingredients {
			recipe.Ingredients[i].RecipeID = recipe.ID
		}
		if err := r.db.Create(&recipe.Ingredients).Error; err != nil {
			return err
		}
	}
	return nil
}

// FIX: UpdateRecipe now handles ingredient replacement properly
func (r *inventoryRepository) UpdateRecipe(recipe *models.Recipe) error {
	// Update recipe base fields (omitting Ingredients to avoid GORM double-saving association)
	if err := r.db.Omit("created_at", "Ingredients").Save(recipe).Error; err != nil {
		return err
	}

	// Delete old ingredients and recreate
	if err := r.db.Where("recipe_id = ?", recipe.ID).Delete(&models.RecipeIngredient{}).Error; err != nil {
		return err
	}

	if len(recipe.Ingredients) > 0 {
		for i := range recipe.Ingredients {
			recipe.Ingredients[i].RecipeID = recipe.ID
			recipe.Ingredients[i].ID = uuid.Nil // Reset ID so GORM creates new records
		}
		if err := r.db.Create(&recipe.Ingredients).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *inventoryRepository) GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.db.Preload("Ingredients").First(&recipe, "menu_id = ?", menuID).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *inventoryRepository) GetRecipeByMenuIDUnscoped(menuID uuid.UUID) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.db.Unscoped().Preload("Ingredients").First(&recipe, "menu_id = ?", menuID).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *inventoryRepository) GetRecipeByID(id uuid.UUID) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.db.Preload("Ingredients").First(&recipe, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

// GetRecipesByMenuIDs — BATCH lookup: 1 query with IN clause + Preload Ingredients in 1 extra query
func (r *inventoryRepository) GetRecipesByMenuIDs(menuIDs []uuid.UUID) ([]models.Recipe, error) {
	if len(menuIDs) == 0 {
		return []models.Recipe{}, nil
	}
	var recipes []models.Recipe
	// Preload Ingredients uses 1 extra batched query (not N+1)
	err := r.db.Preload("Ingredients").Where("menu_id IN ?", menuIDs).Find(&recipes).Error
	return recipes, err
}

func (r *inventoryRepository) DeleteRecipe(id uuid.UUID) error {
	return r.db.Delete(&models.Recipe{}, "id = ?", id).Error
}

func (r *inventoryRepository) HardDeleteRecipe(id uuid.UUID) error {
	// FIX: Also hard delete associated ingredients
	r.db.Unscoped().Where("recipe_id = ?", id).Delete(&models.RecipeIngredient{})
	return r.db.Unscoped().Delete(&models.Recipe{}, "id = ?", id).Error
}
