package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaxService interface {
	CreateTaxConfig(branchID *uuid.UUID, name string, taxType models.TaxType, percentage float64, effectiveFrom time.Time) (*models.TaxConfig, error)
	UpdateTaxConfig(id string, branchID *uuid.UUID, name string, taxType models.TaxType, percentage float64, isActive bool, effectiveFrom time.Time) (*models.TaxConfig, error)
	GetTaxConfigByID(id string) (*models.TaxConfig, error)
	GetActiveTaxes(branchID *uuid.UUID) ([]models.TaxConfig, error)
	DeleteTaxConfig(id string) error
}

type taxService struct {
	db      *gorm.DB
	taxRepo repository.TaxRepository
}

func NewTaxService(db *gorm.DB, taxRepo repository.TaxRepository) TaxService {
	return &taxService{db: db, taxRepo: taxRepo}
}

// CreateTaxConfig — FIXED: Wrap in transaction
func (s *taxService) CreateTaxConfig(branchID *uuid.UUID, name string, taxType models.TaxType, percentage float64, effectiveFrom time.Time) (*models.TaxConfig, error) {
	if name == "" || percentage < 0 || percentage > 1 {
		return nil, errors.New("invalid tax config")
	}
	tax := &models.TaxConfig{
		BranchID: branchID, Name: name, TaxType: taxType,
		Percentage: percentage, IsActive: true, EffectiveFrom: effectiveFrom,
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

	if err := tx.Create(tax).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return tax, nil
}

// UpdateTaxConfig — FIXED: Wrap in transaction
func (s *taxService) UpdateTaxConfig(id string, branchID *uuid.UUID, name string, taxType models.TaxType, percentage float64, isActive bool, effectiveFrom time.Time) (*models.TaxConfig, error) {
	tax, err := s.taxRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("tax config not found")
	}
	if name != "" {
		tax.Name = name
	}
	if taxType != "" {
		tax.TaxType = taxType
	}
	if percentage >= 0 && percentage <= 1 {
		tax.Percentage = percentage
	}
	tax.IsActive = isActive
	if !effectiveFrom.IsZero() {
		tax.EffectiveFrom = effectiveFrom
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

	if err := tx.Save(tax).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return tax, nil
}

func (s *taxService) GetTaxConfigByID(id string) (*models.TaxConfig, error) {
	return s.taxRepo.FindByID(id)
}

func (s *taxService) GetActiveTaxes(branchID *uuid.UUID) ([]models.TaxConfig, error) {
	return s.taxRepo.GetActiveTaxes(branchID)
}

func (s *taxService) DeleteTaxConfig(id string) error {
	return s.taxRepo.Delete(id)
}
