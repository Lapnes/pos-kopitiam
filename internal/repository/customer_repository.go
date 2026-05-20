package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerRepository interface {
	Create(customer *models.Customer) error
	Update(customer *models.Customer) error
	FindByID(id uuid.UUID) (*models.Customer, error)
	FindByPhone(phone string) (*models.Customer, error)
	FindAll() ([]models.Customer, error)
	Delete(id uuid.UUID) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) CustomerRepository {
	return &customerRepository{db: db}
}

func (r *customerRepository) Create(customer *models.Customer) error {
	return r.db.Create(customer).Error
}

func (r *customerRepository) Update(customer *models.Customer) error {
	return r.db.Save(customer).Error
}

func (r *customerRepository) FindByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.First(&customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindByPhone(phone string) (*models.Customer, error) {
	var customer models.Customer
	err := r.db.Where("phone = ?", phone).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *customerRepository) FindAll() ([]models.Customer, error) {
	var customers []models.Customer
	err := r.db.Find(&customers).Error
	return customers, err
}

func (r *customerRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Customer{}, "id = ?", id).Error
}