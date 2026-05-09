package repository

import (
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.Employee, error)
	FindByPIN(pin string) (*models.Employee, error)
	FindByID(id string) (*models.Employee, error)
	FindAll() ([]models.Employee, error)
	Create(user *models.Employee) error
	Update(user *models.Employee) error
	Delete(id string) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByEmail(email string) (*models.Employee, error) {
	var user models.Employee
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByPIN(pin string) (*models.Employee, error) {
	var user models.Employee
	err := r.db.Where("pin_code = ?", pin).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByID(id string) (*models.Employee, error) {
	var user models.Employee
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindAll() ([]models.Employee, error) {
	var users []models.Employee
	err := r.db.Find(&users).Error
	return users, err
}

func (r *UserRepo) Create(user *models.Employee) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) Update(user *models.Employee) error {
	return r.db.Save(user).Error
}

func (r *UserRepo) Delete(id string) error {
	return r.db.Delete(&models.Employee{}, "id = ?", id).Error
}
