package repository

import (
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"gorm.io/gorm"
)

type DiscountRepository interface {
	Create(discount *models.Discount) error
	Update(discount *models.Discount) error
	FindByID(id string) (*models.Discount, error)
	FindActive() ([]models.Discount, error)
	Delete(id string) error

	CreateApplication(app *models.DiscountApplication) error
}

type discountRepository struct {
	db *gorm.DB
}

func NewDiscountRepository(db *gorm.DB) DiscountRepository {
	return &discountRepository{db: db}
}

func (r *discountRepository) Create(discount *models.Discount) error {
	return r.db.Create(discount).Error
}

func (r *discountRepository) Update(discount *models.Discount) error {
	return r.db.Save(discount).Error
}

func (r *discountRepository) FindByID(id string) (*models.Discount, error) {
	var discount models.Discount
	err := r.db.First(&discount, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &discount, nil
}

func (r *discountRepository) FindActive() ([]models.Discount, error) {
	var discounts []models.Discount
	now := time.Now()
	err := r.db.Where("is_active = ? AND (start_date IS NULL OR start_date <= ?) AND (end_date IS NULL OR end_date >= ?)",
		true, now, now).Find(&discounts).Error
	return discounts, err
}

func (r *discountRepository) Delete(id string) error {
	return r.db.Delete(&models.Discount{}, "id = ?", id).Error
}

func (r *discountRepository) CreateApplication(app *models.DiscountApplication) error {
	return r.db.Create(app).Error
}