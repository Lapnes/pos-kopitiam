package repository

import (
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	Create(reservation *models.Reservation) error
	Update(reservation *models.Reservation) error
	FindByID(id string) (*models.Reservation, error)
	FindByBranchID(branchID uuid.UUID) ([]models.Reservation, error)
	FindByTableAndDate(tableID uuid.UUID, date time.Time) ([]models.Reservation, error)
	Delete(id string) error
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) Create(reservation *models.Reservation) error {
	return r.db.Create(reservation).Error
}

func (r *reservationRepository) Update(reservation *models.Reservation) error {
	return r.db.Save(reservation).Error
}

func (r *reservationRepository) FindByID(id string) (*models.Reservation, error) {
	var reservation models.Reservation
	err := r.db.First(&reservation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &reservation, nil
}

func (r *reservationRepository) FindByBranchID(branchID uuid.UUID) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := r.db.Where("branch_id = ?", branchID).Order("reservation_date asc").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) FindByTableAndDate(tableID uuid.UUID, date time.Time) ([]models.Reservation, error) {
	var reservations []models.Reservation
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)
	err := r.db.Where("table_id = ? AND reservation_date BETWEEN ? AND ?", tableID, startOfDay, endOfDay).
		Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) Delete(id string) error {
	return r.db.Delete(&models.Reservation{}, "id = ?", id).Error
}