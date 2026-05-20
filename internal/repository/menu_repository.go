package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MenuRepository interface {
	GetActiveMenusWithCategories() ([]models.Menu, error)
	GetAllMenusWithCategories() ([]models.Menu, error)
	Create(menu *models.Menu) error
	Update(menu *models.Menu) error
	FindByID(id string) (*models.Menu, error)
	FindByIDs(ids []uuid.UUID) ([]models.Menu, error)  // BATCH — avoids N+1
	Delete(id string) error
}

type menuRepository struct {
	db *gorm.DB
}

func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{db: db}
}

// GetActiveMenusWithCategories — Preload is correct for M2M (small set, 2 queries total)
func (r *menuRepository) GetActiveMenusWithCategories() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Preload("Categories").Where("is_active = ?", true).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) GetAllMenusWithCategories() ([]models.Menu, error) {
	var menus []models.Menu
	err := r.db.Preload("Categories").Find(&menus).Error
	return menus, err
}

func (r *menuRepository) Create(menu *models.Menu) error {
	return r.db.Create(menu).Error
}

func (r *menuRepository) Update(menu *models.Menu) error {
	return r.db.Save(menu).Error
}

func (r *menuRepository) FindByID(id string) (*models.Menu, error) {
	var menu models.Menu
	err := r.db.Preload("Categories").First(&menu, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// FindByIDs — BATCH lookup: 1 query with IN clause instead of N queries
func (r *menuRepository) FindByIDs(ids []uuid.UUID) ([]models.Menu, error) {
	if len(ids) == 0 {
		return []models.Menu{}, nil
	}
	var menus []models.Menu
	err := r.db.Where("id IN ?", ids).Find(&menus).Error
	return menus, err
}

func (r *menuRepository) Delete(id string) error {
	return r.db.Delete(&models.Menu{}, "id = ?", id).Error
}
