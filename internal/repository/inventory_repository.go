package repository

import (
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryRepository interface {
	GetDB() *gorm.DB

	// RawMaterial CRUD
	CreateRawMaterial(rm *models.RawMaterial) error
	UpdateRawMaterial(rm *models.RawMaterial) error
	GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error)
	DeleteRawMaterial(id uuid.UUID) error

	// Recipe CRUD
	CreateRecipe(recipe *models.Recipe) error
	UpdateRecipe(recipe *models.Recipe) error
	GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error)
	DeleteRecipe(id uuid.UUID) error
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

// Raw Material Methods
func (r *inventoryRepository) CreateRawMaterial(rm *models.RawMaterial) error {
	return r.db.Create(rm).Error
}

func (r *inventoryRepository) UpdateRawMaterial(rm *models.RawMaterial) error {
	return r.db.Save(rm).Error
}

func (r *inventoryRepository) GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error) {
	var rm models.RawMaterial
	err := r.db.First(&rm, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *inventoryRepository) DeleteRawMaterial(id uuid.UUID) error {
	return r.db.Delete(&models.RawMaterial{}, "id = ?", id).Error
}

// Recipe Methods
func (r *inventoryRepository) CreateRecipe(recipe *models.Recipe) error {
	return r.db.Create(recipe).Error
}

func (r *inventoryRepository) UpdateRecipe(recipe *models.Recipe) error {
	return r.db.Save(recipe).Error
}

func (r *inventoryRepository) GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error) {
	var recipe models.Recipe
	err := r.db.Preload("Ingredients").First(&recipe, "menu_id = ?", menuID).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (r *inventoryRepository) DeleteRecipe(id uuid.UUID) error {
	return r.db.Delete(&models.Recipe{}, "id = ?", id).Error
}
