package service

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InventoryService interface {
	CreateRawMaterial(rm *models.RawMaterial) error
	UpdateRawMaterial(rm *models.RawMaterial) error
	GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error)
	GetAllRawMaterials() ([]models.RawMaterial, error)
	DeleteRawMaterial(id uuid.UUID) error

	UpsertRecipe(recipe *models.Recipe) (*models.Recipe, error)
	UpdateRecipe(recipe *models.Recipe) (*models.Recipe, error)
	GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error)
	DeleteRecipe(id uuid.UUID) error
	GetStockEstimation() ([]map[string]interface{}, error)
}

type inventoryService struct {
	inventoryRepo repository.InventoryRepository
	db            *gorm.DB
}

func NewInventoryService(db *gorm.DB, inventoryRepo repository.InventoryRepository) InventoryService {
	return &inventoryService{db: db, inventoryRepo: inventoryRepo}
}

func (s *inventoryService) CreateRawMaterial(rm *models.RawMaterial) error {
	if rm.Name == "" || rm.Unit == "" {
		return errors.New("name and unit are required")
	}
	return s.inventoryRepo.CreateRawMaterial(rm)
}

func (s *inventoryService) UpdateRawMaterial(rm *models.RawMaterial) error {
	return s.inventoryRepo.UpdateRawMaterial(rm)
}

func (s *inventoryService) GetRawMaterial(id uuid.UUID) (*models.RawMaterial, error) {
	return s.inventoryRepo.GetRawMaterial(id)
}

func (s *inventoryService) GetAllRawMaterials() ([]models.RawMaterial, error) {
	return s.inventoryRepo.GetAllRawMaterials()
}

func (s *inventoryService) DeleteRawMaterial(id uuid.UUID) error {
	return s.inventoryRepo.DeleteRawMaterial(id)
}

// UpsertRecipe — FIXED: Transaction lifecycle bug
func (s *inventoryService) UpsertRecipe(recipe *models.Recipe) (*models.Recipe, error) {
	if len(recipe.Ingredients) == 0 {
		return nil, errors.New("recipe must have at least one ingredient")
	}

	// === TRANSACTION START ===
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	success := false
	defer func() {
		if !success {
			tx.Rollback()
		}
	}()

	// Create tx-scoped repository
	txRepo := repository.NewInventoryRepository(tx)

	existing, err := txRepo.GetRecipeByMenuIDUnscoped(recipe.MenuID)
	if err == nil && existing != nil {
		// Recipe exists - handle update or recreate after soft delete
		if existing.DeletedAt.Valid {
			// Hard delete the soft-deleted recipe first
			if err := txRepo.HardDeleteRecipe(existing.ID); err != nil {
				return nil, fmt.Errorf("failed to clean up deleted recipe: %w", err)
			}
			// Create new recipe with ingredients
			if err := txRepo.CreateRecipe(recipe); err != nil {
				return nil, fmt.Errorf("failed to create recipe after hard delete: %w", err)
			}
			success = true
			tx.Commit()
			// FIX: Gunakan repo utama (bukan txRepo) untuk query setelah commit
			return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
		}
		// Update existing recipe
		recipe.ID = existing.ID
		if err := txRepo.UpdateRecipe(recipe); err != nil {
			return nil, fmt.Errorf("failed to update existing recipe: %w", err)
		}
		success = true
		tx.Commit()
		// FIX: Gunakan repo utama untuk query setelah commit
		return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
	}

	// No existing recipe - create new
	if err := txRepo.CreateRecipe(recipe); err != nil {
		if isDuplicateKeyError(err) {
			// Race condition: recipe was created between our check and insert
			existing, findErr := txRepo.GetRecipeByMenuIDUnscoped(recipe.MenuID)
			if findErr == nil && existing != nil {
				if existing.DeletedAt.Valid {
					txRepo.HardDeleteRecipe(existing.ID)
					if retryErr := txRepo.CreateRecipe(recipe); retryErr != nil {
						return nil, fmt.Errorf("failed after hard delete: %w", retryErr)
					}
					success = true
					tx.Commit()
					return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
				}
				recipe.ID = existing.ID
				if updErr := txRepo.UpdateRecipe(recipe); updErr != nil {
					return nil, fmt.Errorf("failed to update: %w", updErr)
				}
				success = true
				tx.Commit()
				return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
			}
		}
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	success = true
	tx.Commit()
	// FIX: Gunakan repo utama untuk query setelah commit
	return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
}

func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "1062") || strings.Contains(errStr, "Duplicate entry") || strings.Contains(errStr, "duplicate key")
}

func (s *inventoryService) UpdateRecipe(recipe *models.Recipe) (*models.Recipe, error) {
	if err := s.inventoryRepo.UpdateRecipe(recipe); err != nil {
		return nil, err
	}
	return s.inventoryRepo.GetRecipeByMenuID(recipe.MenuID)
}

func (s *inventoryService) GetRecipeByMenuID(menuID uuid.UUID) (*models.Recipe, error) {
	return s.inventoryRepo.GetRecipeByMenuID(menuID)
}

func (s *inventoryService) DeleteRecipe(id uuid.UUID) error {
	return s.inventoryRepo.DeleteRecipe(id)
}

func (s *inventoryService) GetStockEstimation() ([]map[string]interface{}, error) {
	type StockEstimationResult struct {
		RawMaterialID string  `json:"raw_material_id"`
		MaterialName  string  `json:"material_name"`
		Unit          string  `json:"unit"`
		CurrentStock  float64 `json:"current_stock"`
		MinStockLevel float64 `json:"min_stock_level"`
		DailyRequired float64 `json:"daily_required"`
	}

	var results []StockEstimationResult
	err := s.db.Raw(`
		SELECT 
			rm.id as raw_material_id,
			rm.name as material_name,
			rm.unit as unit,
			rm.current_stock as current_stock,
			rm.min_stock_level as min_stock_level,
			COALESCE(SUM(ri.quantity * m.daily_stock), 0) as daily_required
		FROM raw_materials rm
		LEFT JOIN recipe_ingredients ri ON ri.raw_material_id = rm.id AND ri.deleted_at IS NULL
		LEFT JOIN recipes r ON r.id = ri.recipe_id AND r.deleted_at IS NULL
		LEFT JOIN menus m ON m.id = r.menu_id AND m.is_active = 1 AND m.deleted_at IS NULL
		WHERE rm.deleted_at IS NULL
		GROUP BY rm.id, rm.name, rm.unit, rm.current_stock, rm.min_stock_level
	`).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	output := make([]map[string]interface{}, len(results))
	for i, r := range results {
		sufficiency := "sufficient"
		shortage := float64(0)
		if r.CurrentStock < r.DailyRequired {
			sufficiency = "insufficient"
			shortage = r.DailyRequired - r.CurrentStock
		}

		output[i] = map[string]interface{}{
			"raw_material_id": r.RawMaterialID,
			"material_name":   r.MaterialName,
			"unit":            r.Unit,
			"current_stock":   r.CurrentStock,
			"min_stock_level": r.MinStockLevel,
			"daily_required":  r.DailyRequired,
			"sufficiency":     sufficiency,
			"shortage":        shortage,
		}
	}

	return output, nil
}
