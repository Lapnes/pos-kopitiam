package service

import (
	"errors"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerService interface {
	CreateCustomer(name, phone, email string, loyaltyPoints int) (*models.Customer, error)
	UpdateCustomer(id string, name, phone, email string, loyaltyPoints int) (*models.Customer, error)
	GetCustomerByID(id string) (*models.Customer, error)
	GetCustomerByPhone(phone string) (*models.Customer, error)
	GetAllCustomers() ([]models.Customer, error)
	DeleteCustomer(id string) error
}

type customerService struct {
	db           *gorm.DB
	customerRepo repository.CustomerRepository
}

func NewCustomerService(db *gorm.DB, customerRepo repository.CustomerRepository) CustomerService {
	return &customerService{db: db, customerRepo: customerRepo}
}

// CreateCustomer — FIXED: Wrap in transaction
func (s *customerService) CreateCustomer(name, phone, email string, loyaltyPoints int) (*models.Customer, error) {
	if name == "" || phone == "" {
		return nil, errors.New("name and phone are required")
	}
	customer := &models.Customer{
		Name:          name,
		Phone:         phone,
		Email:         email,
		LoyaltyPoints: loyaltyPoints,
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

	if err := tx.Create(customer).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return customer, nil
}

// UpdateCustomer — FIXED: Wrap in transaction
func (s *customerService) UpdateCustomer(id string, name, phone, email string, loyaltyPoints int) (*models.Customer, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	customer, err := s.customerRepo.FindByID(uid)
	if err != nil {
		return nil, errors.New("customer not found")
	}
	if name != "" { customer.Name = name }
	if phone != "" { customer.Phone = phone }
	if email != "" { customer.Email = email }
	if loyaltyPoints >= 0 { customer.LoyaltyPoints = loyaltyPoints }

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

	if err := tx.Save(customer).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return customer, nil
}

func (s *customerService) GetCustomerByID(id string) (*models.Customer, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid id format")
	}
	return s.customerRepo.FindByID(uid)
}

func (s *customerService) GetCustomerByPhone(phone string) (*models.Customer, error) {
	// Assume repository has FindByPhone or we can just query it directly
	var customer models.Customer
	err := s.db.Where("phone = ?", phone).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (s *customerService) GetAllCustomers() ([]models.Customer, error) {
	return s.customerRepo.FindAll()
}

func (s *customerService) DeleteCustomer(id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid id format")
	}
	return s.customerRepo.Delete(uid)
}
