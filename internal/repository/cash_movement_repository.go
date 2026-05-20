package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CashMovementRepository interface {
	Create(cm *models.CashMovement) error
	FindByShiftID(shiftID uuid.UUID) ([]models.CashMovement, error)
}

type cashMovementRepository struct {
	db *gorm.DB
}

func NewCashMovementRepository(db *gorm.DB) CashMovementRepository {
	return &cashMovementRepository{db: db}
}

func (r *cashMovementRepository) Create(cm *models.CashMovement) error {
	return r.db.Create(cm).Error
}

func (r *cashMovementRepository) FindByShiftID(shiftID uuid.UUID) ([]models.CashMovement, error) {
	var movements []models.CashMovement
	err := r.db.Where("shift_id = ?", shiftID).Find(&movements).Error
	return movements, err
}