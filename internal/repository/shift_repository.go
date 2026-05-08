package repository

import (
	"errors"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftRepository interface {
	Create(shift *models.Shift) error
	GetActiveShift(userID uuid.UUID) (*models.Shift, error)
	Update(shift *models.Shift) error
	GetSalesForShift(shiftID uuid.UUID) (totalSales float64, totalTransactions int, err error)
}

type shiftRepository struct {
	db *gorm.DB
}

func NewShiftRepository(db *gorm.DB) ShiftRepository {
	return &shiftRepository{db: db}
}

func (r *shiftRepository) Create(shift *models.Shift) error {
	return r.db.Create(shift).Error
}

func (r *shiftRepository) GetActiveShift(userID uuid.UUID) (*models.Shift, error) {
	var shift models.Shift
	err := r.db.Where("opened_by = ? AND status = ?", userID, "open").First(&shift).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active shift found")
		}
		return nil, err
	}
	return &shift, nil
}

func (r *shiftRepository) Update(shift *models.Shift) error {
	return r.db.Save(shift).Error
}

func (r *shiftRepository) GetSalesForShift(shiftID uuid.UUID) (totalSales float64, totalTransactions int, err error) {
	// Calculate sum of order totals and count of orders for this shift
	var result struct {
		Total float64
		Count int
	}

	err = r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(total), 0) as total, COUNT(id) as count").
		Where("shift_id = ? AND status = ? AND deleted_at IS NULL", shiftID, models.OrderPaid).
		Scan(&result).Error

	if err != nil {
		return 0, 0, err
	}

	return result.Total, result.Count, nil
}
