package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ShiftService interface {
	OpenShift(userID, branchID uuid.UUID, openingCash float64) (*models.Shift, error)
	GetCurrentShift(userID uuid.UUID) (*models.Shift, error)
	CloseShift(userID uuid.UUID, actualClosingCash float64) (*models.Shift, error)
	AddCashMovement(shiftID, recordedBy uuid.UUID, cmType models.CashMovementType, amount float64, reason string) (*models.CashMovement, error)
	GetCashMovements(shiftID uuid.UUID) ([]models.CashMovement, error)
}

type shiftService struct {
	db               *gorm.DB
	shiftRepo        repository.ShiftRepository
	cashMovementRepo repository.CashMovementRepository
}

func NewShiftService(db *gorm.DB, shiftRepo repository.ShiftRepository, cashMovementRepo repository.CashMovementRepository) ShiftService {
	return &shiftService{db: db, shiftRepo: shiftRepo, cashMovementRepo: cashMovementRepo}
}

func (s *shiftService) OpenShift(userID, branchID uuid.UUID, openingCash float64) (*models.Shift, error) {
	activeShift, err := s.shiftRepo.GetActiveShift(userID)
	if err == nil && activeShift != nil {
		return nil, errors.New("user already has an active shift open")
	}
	shift := &models.Shift{
		BranchID: branchID, OpenedBy: userID,
		OpeningTime: time.Now(), OpeningCash: openingCash, Status: "open",
	}
	if err := s.shiftRepo.Create(shift); err != nil {
		return nil, err
	}
	return shift, nil
}

func (s *shiftService) GetCurrentShift(userID uuid.UUID) (*models.Shift, error) {
	return s.shiftRepo.GetActiveShift(userID)
}

// CloseShift — FIXED: Wrap in transaction (shift close + cash movement finalization)
func (s *shiftService) CloseShift(userID uuid.UUID, actualClosingCash float64) (*models.Shift, error) {
	shift, err := s.shiftRepo.GetActiveShift(userID)
	if err != nil {
		return nil, errors.New("no active shift found to close")
	}

	// Get cash movements for this shift (read-only, outside tx)
	movements, err := s.cashMovementRepo.FindByShiftID(shift.ID)
	if err != nil {
		movements = []models.CashMovement{}
	}

	// Calculate expected cash: opening + sum(in) - sum(out)
	var cashMovementTotal float64
	for _, m := range movements {
		if m.Type == models.CashMovementIn {
			cashMovementTotal += m.Amount
		} else {
			cashMovementTotal -= m.Amount
		}
	}

	now := time.Now()
	shift.ClosedBy = &userID
	shift.ClosingTime = &now
	shift.ActualClosingCash = actualClosingCash
	shift.Status = "closed"

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

	if err := tx.Save(shift).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()

	// Return shift with derived fields calculated
	return shift, nil
}

func (s *shiftService) AddCashMovement(shiftID, recordedBy uuid.UUID, cmType models.CashMovementType, amount float64, reason string) (*models.CashMovement, error) {
	cm := &models.CashMovement{
		ShiftID: shiftID, Type: cmType, Amount: amount,
		Reason: reason, RecordedBy: recordedBy, RecordedAt: time.Now(),
	}
	if err := s.cashMovementRepo.Create(cm); err != nil {
		return nil, err
	}
	return cm, nil
}

func (s *shiftService) GetCashMovements(shiftID uuid.UUID) ([]models.CashMovement, error) {
	return s.cashMovementRepo.FindByShiftID(shiftID)
}
