package service

import (
	"errors"
	"fmt"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/google/uuid"
)

type InventoryService interface {
	CreateRawMaterial(rm *models.RawMaterial) error
	UpdateRawMaterial(rm *models.RawMaterial) error
	GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error)
	DeleteRawMaterial(id uuid.UUID) error

	// FIX: Changed CreateRecipe to UpsertRecipe to clarify intent
	UpsertRecipe(recipe *models.Recipe) (*models.Recipe, error)
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

// UpsertRecipe creates a new recipe OR updates existing one if menu_id already exists.
// This handles the case where a recipe already exists for the menu (from seeding or previous creation).
// Returns the created or updated recipe with its ID populated.
func (s *inventoryService) UpsertRecipe(recipe *models.Recipe) (*models.Recipe, error) {
	if len(recipe.Ingredients) == 0 {
		return nil, errors.New("recipe must have at least one ingredient")
	}

	// Try to find existing recipe (including soft-deleted) for this menu
	existing, err := s.inventoryRepo.GetRecipeByMenuIDUnscoped(recipe.MenuID)
	if err == nil && existing != nil {
		if existing.DeletedAt.Valid {
			// Hard delete the old soft-deleted recipe so we can recreate it
			if err := s.inventoryRepo.HardDeleteRecipe(existing.ID); err != nil {
				return nil, fmt.Errorf("failed to clean up deleted recipe: %w", err)
			}
			// Now create the new one
			if err := s.inventoryRepo.CreateRecipe(recipe); err != nil {
				return nil, fmt.Errorf("failed to create recipe after hard delete: %w", err)
			}
			return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
		}

		// Recipe exists and is active - update it (true upsert behavior)
		recipe.ID = existing.ID
		// Preserve created_at by using Omit
		if err := s.inventoryRepo.UpdateRecipe(recipe); err != nil {
			return nil, fmt.Errorf("failed to update existing recipe: %w", err)
		}
		// Reload with associations
		return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
	}

	// No existing recipe found - create new
	if err := s.inventoryRepo.CreateRecipe(recipe); err != nil {
		// Check if error is duplicate key (race condition or soft-deleted record)
		if isDuplicateKeyError(err) {
			// Try to find and update the existing record
			existing, findErr := s.inventoryRepo.GetRecipeByMenuIDUnscoped(recipe.MenuID)
			if findErr == nil && existing != nil {
				if existing.DeletedAt.Valid {
					s.inventoryRepo.HardDeleteRecipe(existing.ID)
					if retryErr := s.inventoryRepo.CreateRecipe(recipe); retryErr != nil {
						return nil, fmt.Errorf("failed to create recipe after duplicate key hard delete: %w", retryErr)
					}
					return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
				}
				
				recipe.ID = existing.ID
				if updErr := s.inventoryRepo.UpdateRecipe(recipe); updErr != nil {
					return nil, fmt.Errorf("failed to update recipe after duplicate key: %w", updErr)
				}
				return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
			}
		}
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	// Reload with associations to ensure ID and relations are populated
	return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return contains(errStr, "1062") || contains(errStr, "Duplicate entry") || contains(errStr, "duplicate key")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsAt(s, substr, 0))
}

func containsAt(s, substr string, start int) bool {
	for i := start; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
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
