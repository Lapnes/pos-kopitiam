package repository

import (
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id string) (*models.Order, error)
	FindByOrderNumber(orderNumber string) (*models.Order, error)
	FindAll() ([]models.Order, error)
	Update(order *models.Order) error
}

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepository {
	return &OrderRepo{db: db}
}

func (r *OrderRepo) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *OrderRepo) FindByID(id string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderDetails").Preload("Payments").First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) FindByOrderNumber(orderNumber string) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderDetails").Preload("Payments").First(&order, "order_number = ?", orderNumber).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepo) FindAll() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("OrderDetails").Preload("Payments").Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepo) Update(order *models.Order) error {
	return r.db.Save(order).Error
}
