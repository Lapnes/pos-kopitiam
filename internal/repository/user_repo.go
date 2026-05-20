package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(email string) (*models.User, error)
	FindByPIN(pin string) (*models.User, error) // tanpa branchID (backward compat)
	FindByID(id string) (*models.User, error)
	FindAll() ([]models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FIX: Gunakan CheckPINHash (cost=10) = 5× lebih cepat
func (r *userRepo) FindByPIN(pin string) (*models.User, error) {
	var users []models.User
	// Hanya fetch active users dengan PIN set
	err := r.db.Where("is_active = ? AND pin != ?", true, "").Find(&users).Error
	if err != nil {
		return nil, err
	}

	// Compare PIN hashes — CheckPINHash menggunakan cost=10 = ~50ms
	for i := range users {
		if utils.CheckPINHash(pin, users[i].PIN) {
			return &users[i], nil
		}
	}

	return nil, gorm.ErrRecordNotFound
}

func (r *userRepo) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepo) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepo) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepo) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepo) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}
