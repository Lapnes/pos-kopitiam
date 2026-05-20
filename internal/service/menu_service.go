package service

import (
	"errors"
	"fmt"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"gorm.io/gorm"
)

type MenuService interface {
	CreateMenu(req models.Menu) (*models.Menu, error)
	UpdateMenu(id string, req models.Menu) (*models.Menu, error)
	GetMenuByID(id string) (*models.Menu, error)
	GetActiveMenus() ([]models.Menu, error)
	GetAllMenus() ([]models.Menu, error)
	DeleteMenu(id string) error
}

type menuService struct {
	db       *gorm.DB
	menuRepo repository.MenuRepository
}

func NewMenuService(db *gorm.DB, menuRepo repository.MenuRepository) MenuService {
	return &menuService{db: db, menuRepo: menuRepo}
}

// CreateMenu — FIXED: Wrap in transaction for atomic create + category association
func (s *menuService) CreateMenu(req models.Menu) (*models.Menu, error) {
	if req.Name == "" || req.Price <= 0 {
		return nil, errors.New("invalid menu data")
	}
	menu := &models.Menu{
		Name:          req.Name,
		Price:         req.Price,
		IsActive:      true,
		IsRecipeBased: req.IsRecipeBased,
		DailyStock:    req.DailyStock,
		Station:       req.Station,
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

	if err := tx.Create(menu).Error; err != nil {
		return nil, err
	}

	// Handle category associations inside transaction
	if len(req.Categories) > 0 {
		if err := tx.Model(menu).Association("Categories").Append(req.Categories); err != nil {
			return nil, fmt.Errorf("failed to associate categories: %w", err)
		}
	}

	success = true
	tx.Commit()
	return menu, nil
}

// UpdateMenu — FIXED: Wrap in transaction for atomic update + category sync
func (s *menuService) UpdateMenu(id string, req models.Menu) (*models.Menu, error) {
	menu, err := s.menuRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("menu not found")
	}
	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Price > 0 {
		menu.Price = req.Price
	}
	if req.Station != "" {
		menu.Station = req.Station
	}
	if req.DailyStock >= 0 {
		menu.DailyStock = req.DailyStock
	}
	menu.IsRecipeBased = req.IsRecipeBased
	menu.IsActive = req.IsActive

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

	if err := tx.Save(menu).Error; err != nil {
		return nil, err
	}

	// Sync categories inside transaction
	if len(req.Categories) > 0 {
		if err := tx.Model(menu).Association("Categories").Replace(req.Categories); err != nil {
			return nil, fmt.Errorf("failed to sync categories: %w", err)
		}
	}

	success = true
	tx.Commit()
	return menu, nil
}

func (s *menuService) GetMenuByID(id string) (*models.Menu, error) {
	return s.menuRepo.FindByID(id)
}

func (s *menuService) GetActiveMenus() ([]models.Menu, error) {
	return s.menuRepo.GetActiveMenusWithCategories()
}

func (s *menuService) GetAllMenus() ([]models.Menu, error) {
	return s.menuRepo.GetAllMenusWithCategories()
}

func (s *menuService) DeleteMenu(id string) error {
	return s.menuRepo.Delete(id)
}
