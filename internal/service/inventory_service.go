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

	CreateRecipe(recipe *models.Recipe) error
	UpdateRecipe(recipe *models.Recipe) error
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

func (s *inventoryService) CreateRecipe(recipe *models.Recipe) error {
	if len(recipe.Ingredients) == 0 {
		return errors.New("recipe must have at least one ingredient")
	}
	return s.inventoryRepo.CreateRecipe(recipe)
}

func (s *inventoryService) UpdateRecipe(recipe *models.Recipe) error {
	return s.inventoryRepo.UpdateRecipe(recipe)
}

func (s *inventoryService) GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error) {
	return s.inventoryRepo.GetRecipeByMenuID(menuID)
}

func (s *inventoryService) DeleteRecipe(id uuid.UUID) error {
	return s.inventoryRepo.DeleteRecipe(id)
}
