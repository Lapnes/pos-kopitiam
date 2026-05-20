package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationService interface {
	CreateReservation(branchID, tableID uuid.UUID, customerID *uuid.UUID, customerName, customerPhone string, reservationDate time.Time, guestCount int, depositAmount float64) (*models.Reservation, error)
	UpdateReservation(id string, customerName, customerPhone string, reservationDate time.Time, guestCount int, depositAmount float64, status models.ReservationStatus) (*models.Reservation, error)
	GetReservationByID(id string) (*models.Reservation, error)
	GetReservationsByBranch(branchID uuid.UUID) ([]models.Reservation, error)
	GetReservationsByTableAndDate(tableID uuid.UUID, date time.Time) ([]models.Reservation, error)
	DeleteReservation(id string) error
}

type reservationService struct {
	db              *gorm.DB
	reservationRepo repository.ReservationRepository
}

func NewReservationService(db *gorm.DB, reservationRepo repository.ReservationRepository) ReservationService {
	return &reservationService{db: db, reservationRepo: reservationRepo}
}

// CreateReservation — FIXED: Wrap in transaction
func (s *reservationService) CreateReservation(branchID, tableID uuid.UUID, customerID *uuid.UUID, customerName, customerPhone string, reservationDate time.Time, guestCount int, depositAmount float64) (*models.Reservation, error) {
	if customerName == "" || customerPhone == "" || guestCount <= 0 {
		return nil, errors.New("invalid reservation data")
	}
	reservation := &models.Reservation{
		BranchID:        branchID,
		TableID:         tableID,
		CustomerID:      customerID,
		CustomerName:    customerName,
		CustomerPhone:   customerPhone,
		ReservationDate: reservationDate,
		GuestCount:      guestCount,
		DepositAmount:   depositAmount,
		Status:          models.ReservationConfirmed,
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

	if err := tx.Create(reservation).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return reservation, nil
}

// UpdateReservation — FIXED: Wrap in transaction
func (s *reservationService) UpdateReservation(id string, customerName, customerPhone string, reservationDate time.Time, guestCount int, depositAmount float64, status models.ReservationStatus) (*models.Reservation, error) {
	reservation, err := s.reservationRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("reservation not found")
	}
	if customerName != "" {
		reservation.CustomerName = customerName
	}
	if customerPhone != "" {
		reservation.CustomerPhone = customerPhone
	}
	reservation.ReservationDate = reservationDate
	reservation.GuestCount = guestCount
	reservation.DepositAmount = depositAmount
	reservation.Status = status

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

	if err := tx.Save(reservation).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return reservation, nil
}

func (s *reservationService) GetReservationByID(id string) (*models.Reservation, error) {
	return s.reservationRepo.FindByID(id)
}

func (s *reservationService) GetReservationsByBranch(branchID uuid.UUID) ([]models.Reservation, error) {
	return s.reservationRepo.FindByBranchID(branchID)
}

func (s *reservationService) GetReservationsByTableAndDate(tableID uuid.UUID, date time.Time) ([]models.Reservation, error) {
	return s.reservationRepo.FindByTableAndDate(tableID, date)
}

func (s *reservationService) DeleteReservation(id string) error {
	return s.reservationRepo.Delete(id)
}
