package service

import (
	"errors"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockService interface {
	AdjustStock(rawMaterialID uuid.UUID, quantityAfter float64, reason string, adjustedBy uuid.UUID) (*models.StockAdjustment, error)
	GetAdjustmentsByRawMaterial(rawMaterialID uuid.UUID) ([]models.StockAdjustment, error)
}

type stockService struct {
	stockRepo     repository.StockRepository
	inventoryRepo repository.InventoryRepository
	db            *gorm.DB
}

func NewStockService(db *gorm.DB, stockRepo repository.StockRepository, inventoryRepo repository.InventoryRepository) StockService {
	return &stockService{db: db, stockRepo: stockRepo, inventoryRepo: inventoryRepo}
}

// AdjustStock — FIXED: Wrap in transaction (raw material update + adjustment record must be atomic)
func (s *stockService) AdjustStock(rawMaterialID uuid.UUID, quantityAfter float64, reason string, adjustedBy uuid.UUID) (*models.StockAdjustment, error) {
	rm, err := s.inventoryRepo.GetRawMaterial(rawMaterialID)
	if err != nil {
		return nil, errors.New("raw material not found")
	}

	adjustment := &models.StockAdjustment{
		RawMaterialID:  rawMaterialID,
		QuantityBefore: rm.CurrentStock,
		QuantityAfter:  quantityAfter,
		Reason:         reason,
		AdjustedBy:     adjustedBy,
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

	// Update raw material stock inside transaction
	rm.CurrentStock = quantityAfter
	if err := tx.Save(rm).Error; err != nil {
		return nil, err
	}

	// Create adjustment record inside same transaction
	if err := tx.Create(adjustment).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return adjustment, nil
}

func (s *stockService) GetAdjustmentsByRawMaterial(rawMaterialID uuid.UUID) ([]models.StockAdjustment, error) {
	return s.stockRepo.FindByRawMaterialID(rawMaterialID)
}
