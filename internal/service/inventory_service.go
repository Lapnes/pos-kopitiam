package service

import (
	"errors"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/google/uuid"
)

type InventoryService interface {
	CreateRawMaterial(rm *models.RawMaterial) error
	UpdateRawMaterial(rm *models.RawMaterial) error
	GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error)
	DeleteRawMaterial(id uuid.UUID) error

	CreateRecipe(recipe *models.Recipe) (*models.Recipe, error)
	UpdateRecipe(recipe *models.Recipe) (*models.Recipe, error)
	GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error)
	DeleteRecipe(id uuid.UUID) error
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
}

func NewInventoryService(inventoryRepo repository.InventoryRepository) InventoryService {
	return &inventoryService{inventoryRepo: inventoryRepo}
}

func (s *inventoryService) CreateRawMaterial(rm *models.RawMaterial) error {
	if rm.Name == "" {
		return errors.New("raw material name is required")
	}
	return s.inventoryRepo.CreateRawMaterial(rm)
}

func (s *inventoryService) UpdateRawMaterial(rm *models.RawMaterial) error {
	return s.inventoryRepo.UpdateRawMaterial(rm)
}

func (s *inventoryService) GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error) {
	return s.inventoryRepo.GetRawMaterial(id)
}

func (s *inventoryService) DeleteRawMaterial(id uuid.UUID) error {
	return s.inventoryRepo.DeleteRawMaterial(id)
}

// CreateRecipe creates a new recipe or updates existing one if menu_id already exists (upsert).
// Returns the created or updated recipe with its ID populated.
func (s *inventoryService) CreateRecipe(recipe *models.Recipe) (*models.Recipe, error) {
	if len(recipe.Ingredients) == 0 {
		return nil, errors.New("recipe must have at least one ingredient")
	}

	db := s.inventoryRepo.GetDB()
	if db == nil {
		return nil, errors.New("database connection not available")
	}

	// Try to find existing recipe first
	existing, err := s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
	if err == nil && existing != nil {
		// Recipe exists, update it instead
		recipe.ID = existing.ID
		if err := s.inventoryRepo.UpdateRecipe(recipe); err != nil {
			return nil, err
		}
		// Reload with associations
		return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
	}

	// No existing recipe, create new
	if err := s.inventoryRepo.CreateRecipe(recipe); err != nil {
		return nil, err
	}

	// Reload with associations to ensure ID and relations are populated
	return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
}

func (s *inventoryService) UpdateRecipe(recipe *models.Recipe) (*models.Recipe, error) {
	if err := s.inventoryRepo.UpdateRecipe(recipe); err != nil {
		return nil, err
	}
	// Reload with associations
	return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
}

func (s *inventoryService) GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error) {
	return s.inventoryRepo.GetRecipeByMenuID(menuID)
}

func (s *inventoryService) DeleteRecipe(id uuid.UUID) error {
	return s.inventoryRepo.DeleteRecipe(id)
}
