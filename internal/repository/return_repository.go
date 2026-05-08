package repository

import (
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/gorm"
)

type ReturnRepository interface {
	Create(orderReturn *models.OrderReturn) error
}

type returnRepository struct {
	db *gorm.DB
}

func NewReturnRepository(db *gorm.DB) ReturnRepository {
	return &returnRepository{db: db}
}

func (r *returnRepository) Create(orderReturn *models.OrderReturn) error {
	return r.db.Create(orderReturn).Error
}
