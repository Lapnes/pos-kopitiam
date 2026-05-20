package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaxRepository interface {
	Create(tax *models.TaxConfig) error
	Update(tax *models.TaxConfig) error
	FindByID(id string) (*models.TaxConfig, error)
	GetActiveTaxes(branchID *uuid.UUID) ([]models.TaxConfig, error)
	Delete(id string) error
}

type taxRepository struct {
	db *gorm.DB
}

func NewTaxRepository(db *gorm.DB) TaxRepository {
	return &taxRepository{db: db}
}

func (r *taxRepository) Create(tax *models.TaxConfig) error {
	return r.db.Create(tax).Error
}

func (r *taxRepository) Update(tax *models.TaxConfig) error {
	return r.db.Save(tax).Error
}

func (r *taxRepository) FindByID(id string) (*models.TaxConfig, error) {
	var tax models.TaxConfig
	err := r.db.First(&tax, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tax, nil
}

func (r *taxRepository) GetActiveTaxes(branchID *uuid.UUID) ([]models.TaxConfig, error) {
	var taxes []models.TaxConfig
	query := r.db.Where("is_active = ? AND effective_from <= NOW()", true)
	if branchID != nil {
		query = query.Where("branch_id = ? OR branch_id IS NULL", *branchID)
	} else {
		query = query.Where("branch_id IS NULL")
	}
	err := query.Find(&taxes).Error
	return taxes, err
}

func (r *taxRepository) Delete(id string) error {
	return r.db.Delete(&models.TaxConfig{}, "id = ?", id).Error
}