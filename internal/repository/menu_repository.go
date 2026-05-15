package repository

import (
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/gorm"
)

type MenuRepository interface {
	GetActiveMenus() ([]models.Menu, error)
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

func (r *menuRepository) GetActiveMenus() ([]models.Menu, error) {
	var menus []models.Menu
	// Preload Categories for M2M relationship
	err := r.db.Preload("Categories").Where("is_active = ?", true).Find(&menus).Error
	return menus, err
}
