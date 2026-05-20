package service

import (
	"errors"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"gorm.io/gorm"
)

type CategoryService interface {
	CreateCategory(name string, sortOrder int) (*models.Category, error)
	UpdateCategory(id string, name string, sortOrder int) (*models.Category, error)
	GetCategoryByID(id string) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	DeleteCategory(id string) error
}

type categoryService struct {
	db           *gorm.DB
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(db *gorm.DB, categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{db: db, categoryRepo: categoryRepo}
}

// CreateCategory — FIXED: Wrap in transaction
func (s *categoryService) CreateCategory(name string, sortOrder int) (*models.Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}
	category := &models.Category{Name: name, SortOrder: sortOrder}

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

	if err := tx.Create(category).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return category, nil
}

// UpdateCategory — FIXED: Wrap in transaction
func (s *categoryService) UpdateCategory(id string, name string, sortOrder int) (*models.Category, error) {
	category, err := s.categoryRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("category not found")
	}
	if name != "" {
		category.Name = name
	}
	category.SortOrder = sortOrder

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

	if err := tx.Save(category).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return category, nil
}

func (s *categoryService) GetCategoryByID(id string) (*models.Category, error) {
	return s.categoryRepo.FindByID(id)
}

func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.FindAll()
}

func (s *categoryService) DeleteCategory(id string) error {
	return s.categoryRepo.Delete(id)
}
