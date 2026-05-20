package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	Create(payment *models.Payment) error
	FindByOrderID(orderID uuid.UUID) ([]models.Payment, error)
	FindByID(id uuid.UUID) (*models.Payment, error)
	Update(payment *models.Payment) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

func (r *paymentRepository) FindByOrderID(orderID uuid.UUID) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Preload("Splits").Where("order_id = ?", orderID).Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) FindByID(id uuid.UUID) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Preload("Splits").First(&payment, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Save(payment).Error
}