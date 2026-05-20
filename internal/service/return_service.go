package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReturnService interface {
	ProcessReturn(orderID, processedBy uuid.UUID, returnAmount float64, reason string) (*models.OrderReturn, error)
}

type returnService struct {
	db           *gorm.DB
	returnRepo   repository.ReturnRepository
	orderRepo    repository.OrderRepository
	orderService OrderService
}

func NewReturnService(db *gorm.DB, returnRepo repository.ReturnRepository, orderRepo repository.OrderRepository, orderService OrderService) ReturnService {
	return &returnService{db: db, returnRepo: returnRepo, orderRepo: orderRepo, orderService: orderService}
}

// ProcessReturn — FIXED: Proper transaction boundary + datetime handling
func (s *returnService) ProcessReturn(orderID, processedBy uuid.UUID, returnAmount float64, reason string) (*models.OrderReturn, error) {
	order, err := s.orderRepo.FindByID(orderID.String())
	if err != nil {
		return nil, errors.New("order not found")
	}

	// Calculate order total dynamically
	var orderTotal float64
	for _, detail := range order.OrderDetails {
		if !detail.IsVoided {
			orderTotal += detail.Subtotal
		}
	}
	// Add tax & service if needed
	orderTotal += order.TaxAmount + order.ServiceCharge - order.DiscountAmount
	if orderTotal < 0 {
		orderTotal = 0
	}

	if returnAmount > orderTotal {
		return nil, errors.New("return amount exceeds order total")
	}

	// Pre-collect menu IDs for batch lookup
	menuIDs := make([]uuid.UUID, 0, len(order.OrderDetails))
	for _, detail := range order.OrderDetails {
		menuIDs = append(menuIDs, detail.MenuID)
	}

	// Batch fetch menus (1 query instead of N)
	menuMap := make(map[uuid.UUID]models.Menu)
	if len(menuIDs) > 0 {
		var menus []models.Menu
		if err := s.db.Where("id IN ?", menuIDs).Find(&menus).Error; err == nil {
			for _, m := range menus {
				menuMap[m.ID] = m
			}
		}
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

	orderReturn := &models.OrderReturn{
		OrderID:        orderID,
		ProcessedBy:    processedBy,
		Reason:         reason,
		OriginalAmount: orderTotal,
		ReturnAmount:   returnAmount,
		ProcessedAt:    time.Now(), // FIX: Set explicit, valid datetime
	}

	// FIX: Hapus manual CreatedAt/UpdatedAt — GORM auto-manage
	if err := tx.Create(orderReturn).Error; err != nil {
		return nil, err
	}

	order.Status = models.OrderRefunded
	order.HasReturns = true
	if err := tx.Save(order).Error; err != nil {
		return nil, err
	}

	// Restore stock using batch-lookup map (no per-item DB queries)
	for _, detail := range order.OrderDetails {
		if menu, ok := menuMap[detail.MenuID]; ok {
			menu.DailyStock += detail.Quantity
			if err := tx.Save(&menu).Error; err != nil {
				return nil, errors.New("failed to restore stock")
			}
		}
	}

	success = true
	tx.Commit()
	return orderReturn, nil
}
