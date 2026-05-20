package service

import (
	"errors"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EmployeeService interface {
	CreateEmployee(name, email, password string, role models.Role, pin string, branchID uuid.UUID) (*models.User, error)
	UpdateEmployee(id string, name, email, password string, role models.Role, pin string, isActive bool) (*models.User, error)
	GetEmployeeByID(id string) (*models.User, error)
	GetAllEmployees() ([]models.User, error)
	DeleteEmployee(id string) error
}

type employeeService struct {
	db       *gorm.DB
	userRepo repository.UserRepository
}

func NewEmployeeService(db *gorm.DB, userRepo repository.UserRepository) EmployeeService {
	return &employeeService{db: db, userRepo: userRepo}
}

// CreateEmployee — FIXED: Hash password and PIN properly, accept password parameter
func (s *employeeService) CreateEmployee(name, email, password string, role models.Role, pin string, branchID uuid.UUID) (*models.User, error) {
	if name == "" || email == "" || pin == "" {
		return nil, errors.New("name, email, and PIN are required")
	}

	// FIX: Use provided password or generate default if empty
	passwordToHash := password
	if passwordToHash == "" {
		passwordToHash = "password123" // Default password
	}

	passwordHash, err := utils.HashPassword(passwordToHash) // password tetap cost=14
	if err != nil {
		return nil, err
	}
	pinHash, err := utils.HashPIN(pin) // ← FIX: PIN cost=10
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		Phone:        "",
		PasswordHash: passwordHash,
		PIN:          pinHash,
		Role:         role,
		BranchID:     branchID,
		IsActive:     true,
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

	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return user, nil
}

// UpdateEmployee — FIXED: Handle password updates properly
func (s *employeeService) UpdateEmployee(id string, name, email, password string, role models.Role, pin string, isActive bool) (*models.User, error) {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}
	if name != "" {
		user.Name = name
	}
	if email != "" {
		user.Email = email
	}
	// FIX: Update password hash if new password provided
	if password != "" {
		newHash, err := utils.HashPassword(password)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = newHash
	}
	// FIX: Update PIN hash if new PIN provided
	if pin != "" {
		newPinHash, err := utils.HashPIN(pin)
		if err != nil {
			return nil, err
		}
		user.PIN = newPinHash
	}
	if role != "" {
		user.Role = role
	}
	user.IsActive = isActive

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

	if err := tx.Save(user).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return user, nil
}

func (s *employeeService) GetEmployeeByID(id string) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *employeeService) GetAllEmployees() ([]models.User, error) {
	return s.userRepo.FindAll()
}

func (s *employeeService) DeleteEmployee(id string) error {
	return s.userRepo.Delete(id)
}
