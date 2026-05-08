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
	// Preload Category if needed, but for now just fetching the menus is enough
	err := r.db.Where("is_active = ?", true).Find(&menus).Error
	return menus, err
}
