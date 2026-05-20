package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StockRepository interface {
	Create(adjustment *models.StockAdjustment) error
	FindByRawMaterialID(rawMaterialID uuid.UUID) ([]models.StockAdjustment, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db: db}
}

func (r *stockRepository) Create(adjustment *models.StockAdjustment) error {
	return r.db.Create(adjustment).Error
}

func (r *stockRepository) FindByRawMaterialID(rawMaterialID uuid.UUID) ([]models.StockAdjustment, error) {
	var adjustments []models.StockAdjustment
	err := r.db.Where("raw_material_id = ?", rawMaterialID).Order("adjusted_at desc").Find(&adjustments).Error
	return adjustments, err
}