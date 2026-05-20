package repository

import (
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"gorm.io/gorm"
)

type ReturnRepository interface {
	Create(orderReturn *models.OrderReturn) error
	FindByOrderID(orderID string) ([]models.OrderReturn, error)
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

func (r *returnRepository) FindByOrderID(orderID string) ([]models.OrderReturn, error) {
	var returns []models.OrderReturn
	err := r.db.Where("order_id = ?", orderID).Find(&returns).Error
	return returns, err
}