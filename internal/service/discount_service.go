package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscountService interface {
	CreateDiscount(name string, discountType models.DiscountType, value, minOrder, maxDiscount float64, startDate, endDate *time.Time) (*models.Discount, error)
	UpdateDiscount(id string, name string, discountType models.DiscountType, value, minOrder, maxDiscount float64, startDate, endDate *time.Time, isActive bool) (*models.Discount, error)
	GetDiscountByID(id string) (*models.Discount, error)
	GetActiveDiscounts() ([]models.Discount, error)
	DeleteDiscount(id string) error
	ApplyDiscount(orderID, discountID uuid.UUID, discountAmount float64) (*models.DiscountApplication, error)
}

type discountService struct {
	db           *gorm.DB
	discountRepo repository.DiscountRepository
}

func NewDiscountService(db *gorm.DB, discountRepo repository.DiscountRepository) DiscountService {
	return &discountService{db: db, discountRepo: discountRepo}
}

func (s *discountService) CreateDiscount(name string, discountType models.DiscountType, value, minOrder, maxDiscount float64, startDate, endDate *time.Time) (*models.Discount, error) {
	if name == "" || value <= 0 {
		return nil, errors.New("invalid discount parameters")
	}
	discount := &models.Discount{
		Name:           name,
		Type:           discountType,
		Value:          value,
		MinOrderAmount: minOrder,
		MaxDiscount:    maxDiscount,
		StartDate:      startDate,
		EndDate:        endDate,
		IsActive:       true,
	}
	if err := s.discountRepo.Create(discount); err != nil {
		return nil, err
	}
	return discount, nil
}

func (s *discountService) UpdateDiscount(id string, name string, discountType models.DiscountType, value, minOrder, maxDiscount float64, startDate, endDate *time.Time, isActive bool) (*models.Discount, error) {
	discount, err := s.discountRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("discount not found")
	}
	if name != "" {
		discount.Name = name
	}
	if discountType != "" {
		discount.Type = discountType
	}
	if value > 0 {
		discount.Value = value
	}
	if minOrder >= 0 {
		discount.MinOrderAmount = minOrder
	}
	if maxDiscount >= 0 {
		discount.MaxDiscount = maxDiscount
	}
	if startDate != nil {
		discount.StartDate = startDate
	}
	if endDate != nil {
		discount.EndDate = endDate
	}
	discount.IsActive = isActive
	if err := s.discountRepo.Update(discount); err != nil {
		return nil, err
	}
	return discount, nil
}

func (s *discountService) GetDiscountByID(id string) (*models.Discount, error) {
	return s.discountRepo.FindByID(id)
}

func (s *discountService) GetActiveDiscounts() ([]models.Discount, error) {
	return s.discountRepo.FindActive()
}

func (s *discountService) DeleteDiscount(id string) error {
	return s.discountRepo.Delete(id)
}

// ApplyDiscount — FIXED: Wrap in transaction for atomicity
func (s *discountService) ApplyDiscount(orderID, discountID uuid.UUID, discountAmount float64) (*models.DiscountApplication, error) {
	app := &models.DiscountApplication{
		OrderID:        orderID,
		DiscountID:     discountID,
		DiscountAmount: discountAmount,
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

	if err := tx.Create(app).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return app, nil
}
