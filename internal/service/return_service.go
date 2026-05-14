package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReturnService interface {
	ProcessReturn(orderID, processedBy uuid.UUID, returnAmount float64, reason string) (*models.OrderReturn, error)
}

type returnService struct {
	db         *gorm.DB
	returnRepo repository.ReturnRepository
	orderRepo  repository.OrderRepository
}

func NewReturnService(db *gorm.DB, returnRepo repository.ReturnRepository, orderRepo repository.OrderRepository) ReturnService {
	return &returnService{
		db:         db,
		returnRepo: returnRepo,
		orderRepo:  orderRepo,
	}
}

func (s *returnService) ProcessReturn(orderID, processedBy uuid.UUID, returnAmount float64, reason string) (*models.OrderReturn, error) {
	order, err := s.orderRepo.FindByID(orderID.String())
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}

	if returnAmount > order.Total {
		return nil, errors.New("return amount exceeds the order total")
	}

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

	orderReturn := &models.OrderReturn{
		OrderID:        orderID,
		ProcessedBy:    processedBy,
		Reason:         reason,
		OriginalAmount: order.Total,
		ReturnAmount:   returnAmount,
	}

	now := time.Now()
	orderReturn.CreatedAt = now
	orderReturn.UpdatedAt = now

	if err := tx.Create(orderReturn).Error; err != nil {
		return nil, errors.New("failed to save return record: " + err.Error())
	}

	order.Status = models.OrderRefunded
	order.HasReturns = true
	if err := tx.Save(order).Error; err != nil {
		return nil, errors.New("failed to update order status: " + err.Error())
	}

	// Restore stock
	for _, detail := range order.OrderDetails {
		var menu models.Menu
		if err := tx.Where("id = ?", detail.MenuID).First(&menu).Error; err == nil {
			menu.DailyStock += detail.Quantity
			tx.Save(&menu)
		}
	}

	success = true
	tx.Commit()

	return orderReturn, nil
}
