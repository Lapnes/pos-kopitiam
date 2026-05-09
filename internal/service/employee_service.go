package service

import (
	"errors"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/google/uuid"
)

type EmployeeService interface {
	GetAllEmployees() ([]models.Employee, error)
	GetEmployeeByID(id string) (*models.Employee, error)
	CreateEmployee(branchID uuid.UUID, name, email, password, pin string, role models.Role, isActive bool) (*models.Employee, error)
	UpdateEmployee(id string, name, email, password, pin string, role models.Role, isActive bool) (*models.Employee, error)
	DeleteEmployee(id string) error
}

type employeeService struct {
	userRepo repository.UserRepository
}

func NewEmployeeService(userRepo repository.UserRepository) EmployeeService {
	return &employeeService{userRepo: userRepo}
}

func (s *employeeService) GetAllEmployees() ([]models.Employee, error) {
	return s.userRepo.FindAll()
}

func (s *employeeService) GetEmployeeByID(id string) (*models.Employee, error) {
	return s.userRepo.FindByID(id)
}

func (s *employeeService) CreateEmployee(branchID uuid.UUID, name, email, password, pin string, role models.Role, isActive bool) (*models.Employee, error) {
	// Hash password and pin
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	hashedPin, err := utils.HashPassword(pin)
	if err != nil {
		return nil, errors.New("failed to hash pin")
	}

	employee := &models.Employee{
		BranchID: branchID,
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		PINCode:  hashedPin,
		Role:     role,
		IsActive: isActive,
	}

	if err := s.userRepo.Create(employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *employeeService) UpdateEmployee(id string, name, email, password, pin string, role models.Role, isActive bool) (*models.Employee, error) {
	employee, err := s.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("employee not found")
	}

	employee.Name = name
	employee.Email = email
	employee.Role = role
	employee.IsActive = isActive

	if password != "" {
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		employee.Password = hashedPassword
	}

	if pin != "" {
		hashedPin, err := utils.HashPassword(pin)
		if err != nil {
			return nil, errors.New("failed to hash pin")
		}
		employee.PINCode = hashedPin
	}

	if err := s.userRepo.Update(employee); err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *employeeService) DeleteEmployee(id string) error {
	return s.userRepo.Delete(id)
}
