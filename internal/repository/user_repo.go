package repository

import (
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.Employee, error)
	FindByPIN(pin string) (*models.Employee, error)
	FindByID(id string) (*models.Employee, error)
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
