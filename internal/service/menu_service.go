package service

import (
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
)

type MenuService interface {
	GetActiveMenus() ([]models.Menu, error)
}

type menuService struct {
	menuRepo repository.MenuRepository
}

func NewMenuService(menuRepo repository.MenuRepository) MenuService {
	return &menuService{menuRepo: menuRepo}
}

func (s *menuService) GetActiveMenus() ([]models.Menu, error) {
	return s.menuRepo.GetActiveMenus()
}
