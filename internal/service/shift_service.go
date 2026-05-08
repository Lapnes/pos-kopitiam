package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/google/uuid"
)

type ShiftService interface {
	OpenShift(userID, branchID uuid.UUID, openingCash float64) (*models.Shift, error)
	GetCurrentShift(userID uuid.UUID) (*models.Shift, error)
	CloseShift(userID uuid.UUID, actualClosingCash float64) (*models.Shift, error)
}

type shiftService struct {
	shiftRepo repository.ShiftRepository
}

func NewShiftService(shiftRepo repository.ShiftRepository) ShiftService {
	return &shiftService{shiftRepo: shiftRepo}
}

func (s *shiftService) OpenShift(userID, branchID uuid.UUID, openingCash float64) (*models.Shift, error) {
	// Check if there is already an active shift for this user
	activeShift, err := s.shiftRepo.GetActiveShift(userID)
	if err == nil && activeShift != nil {
		return nil, errors.New("user already has an active shift open")
	}

	shift := &models.Shift{
		BranchID:    branchID,
		OpenedBy:    userID,
		OpeningTime: time.Now(),
		OpeningCash: openingCash,
		Status:      "open",
	}

	if err := s.shiftRepo.Create(shift); err != nil {
		return nil, err
	}

	return shift, nil
}

func (s *shiftService) GetCurrentShift(userID uuid.UUID) (*models.Shift, error) {
	return s.shiftRepo.GetActiveShift(userID)
}

func (s *shiftService) CloseShift(userID uuid.UUID, actualClosingCash float64) (*models.Shift, error) {
	shift, err := s.shiftRepo.GetActiveShift(userID)
	if err != nil {
		return nil, errors.New("no active shift found to close")
	}

	// Calculate Sales
	totalSales, totalTransactions, err := s.shiftRepo.GetSalesForShift(shift.ID)
	if err != nil {
		return nil, errors.New("failed to calculate shift sales: " + err.Error())
	}

	expectedCash := shift.OpeningCash + totalSales
	variance := actualClosingCash - expectedCash

	now := time.Now()

	shift.ClosedBy = &userID
	shift.ClosingTime = &now
	shift.ActualClosingCash = actualClosingCash
	shift.ExpectedCash = expectedCash
	shift.Variance = variance
	shift.TotalSales = totalSales
	shift.TotalTransactions = totalTransactions
	shift.Status = "closed"

	if err := s.shiftRepo.Update(shift); err != nil {
		return nil, err
	}

	return shift, nil
}
