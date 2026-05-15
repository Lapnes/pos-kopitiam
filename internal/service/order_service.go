package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/dto"
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderService interface {
	CreateOrder(req dto.CreateOrderRequest, userID string, branchID string) (*models.Order, error)
	GetOrders() ([]models.Order, error)
	GetOrder(id string) (*models.Order, error)
	GetOrderByNumber(orderNumber string) (*models.Order, error)
	UpdateOrderNotes(id string, notes string) (*models.Order, error)
	ConfirmOrder(orderID string) (*models.Order, error)
	CancelOrder(orderID string) (*models.Order, error)
	VoidItem(orderID string, orderDetailID string, reason string) (*models.Order, error)
	VoidOrder(orderID string, reason string) (*models.Order, error)
}

type orderService struct {
	db            *gorm.DB
	orderRepo     repository.OrderRepository
	inventoryRepo repository.InventoryRepository
	shiftRepo     repository.ShiftRepository
}

// FIX: Updated constructor to accept shiftRepo
func NewOrderService(db *gorm.DB, orderRepo repository.OrderRepository, inventoryRepo repository.InventoryRepository, shiftRepo repository.ShiftRepository) OrderService {
	return &orderService{
		db:            db,
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		shiftRepo:     shiftRepo,
	}
}

// generateOrderNumber creates a unique order number using crypto-secure random
func generateOrderNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("ORD-%s-%04d", time.Now().Format("20060102"), n)
}

func (s *orderService) CreateOrder(req dto.CreateOrderRequest, userID string, branchID string) (*models.Order, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	orderNumber := generateOrderNumber()

	branchUUID, _ := uuid.Parse(branchID)
	userUUID, _ := uuid.Parse(userID)
	var tableUUID *uuid.UUID
	if req.TableID != "" {
		id, _ := uuid.Parse(req.TableID)
		tableUUID = &id
	}

	var orderDetails []models.OrderDetail
	for _, item := range req.Items {
		itemUUID, _ := uuid.Parse(item.MenuID)
		itemSubtotal := item.Price * float64(item.Quantity)

		var menu models.Menu
		if err := s.db.Where("id = ?", itemUUID).First(&menu).Error; err != nil {
			return nil, fmt.Errorf("menu item %s not found: %w", item.MenuID, err)
		}

		// FIX: Validate menu is active
		if !menu.IsActive {
			return nil, fmt.Errorf("menu item %s is not active", menu.Name)
		}

		orderDetails = append(orderDetails, models.OrderDetail{
			MenuID:    itemUUID,
			MenuName:  item.MenuName,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CostPrice: menu.CostPrice,
			Subtotal:  itemSubtotal,
			Notes:     item.Notes,
			Station:   menu.Station,
		})
	}

	order := &models.Order{
		OrderNumber:    orderNumber,
		BranchID:       branchUUID,
		TableID:        tableUUID,
		EmployeeID:     userUUID,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		OrderType:      models.OrderType(req.OrderType),
		Status:         models.OrderPending,
		DiscountAmount: req.DiscountAmount,
		Notes:          req.Notes,
		OrderDetails:   orderDetails,
	}

	order.CalculateTotal()

	err := s.orderRepo.Create(order)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (s *orderService) GetOrders() ([]models.Order, error) {
	return s.orderRepo.FindAll()
}

func (s *orderService) GetOrder(id string) (*models.Order, error) {
	return s.orderRepo.FindByID(id)
}

func (s *orderService) GetOrderByNumber(orderNumber string) (*models.Order, error) {
	return s.orderRepo.FindByOrderNumber(orderNumber)
}

func (s *orderService) UpdateOrderNotes(id string, notes string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}
	order.Notes = notes
	if err := s.orderRepo.Update(order); err != nil {
		return nil, errors.New("failed to update order: " + err.Error())
	}
	return order, nil
}

func (s *orderService) ConfirmOrder(orderID string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}

	if order.Status != models.OrderPending {
		return nil, errors.New("only pending orders can be confirmed")
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

	for _, detail := range order.OrderDetails {
		var menu models.Menu
		if err := tx.Where("id = ?", detail.MenuID).First(&menu).Error; err != nil {
			return nil, fmt.Errorf("menu item %s not found during confirm: %w", detail.MenuName, err)
		}

		if menu.IsRecipeBased {
			// FIX: Use Unscoped to find soft-deleted recipes too, then check if active
			var recipe models.Recipe
			if err := tx.Unscoped().Preload("Ingredients").Where("menu_id = ?", detail.MenuID).First(&recipe).Error; err != nil {
				return nil, fmt.Errorf("recipe not found for recipe-based menu %s: %w", detail.MenuName, err)
			}

			// If recipe is soft-deleted, we can't use it for stock deduction
			if recipe.DeletedAt.Valid {
				return nil, fmt.Errorf("recipe for menu %s has been deleted", detail.MenuName)
			}

			for _, ingredient := range recipe.Ingredients {
				totalDeduction := ingredient.Quantity * float64(detail.Quantity)

				var rm models.RawMaterial
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", ingredient.RawMaterialID).First(&rm).Error; err != nil {
					return nil, errors.New("raw material not found")
				}

				if rm.CurrentStock < totalDeduction {
					return nil, fmt.Errorf("insufficient stock for raw material: %s", rm.Name)
				}

				rm.CurrentStock -= totalDeduction
				if err := tx.Save(&rm).Error; err != nil {
					return nil, errors.New("failed to update raw material stock")
				}
			}
		} else {
			if menu.DailyStock < detail.Quantity {
				return nil, fmt.Errorf("insufficient daily stock for menu: %s", menu.Name)
			}
			menu.DailyStock -= detail.Quantity
			if menu.DailyStock < 0 {
				menu.DailyStock = 0
			}
			if err := tx.Save(&menu).Error; err != nil {
				return nil, errors.New("failed to update menu daily stock")
			}
		}
	}

	order.Status = models.OrderConfirmed

	// FIX: Set ShiftID from current active shift if available
	if s.shiftRepo != nil {
		shift, err := s.shiftRepo.GetActiveShift(order.EmployeeID)
		if err == nil && shift != nil {
			order.ShiftID = &shift.ID
		}
	}

	if err := tx.Save(order).Error; err != nil {
		return nil, errors.New("failed to update order status")
	}

	success = true
	tx.Commit()

	return order, nil
}

// CancelOrder cancels a pending order (cashier action)
func (s *orderService) CancelOrder(orderID string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}

	if order.Status != models.OrderPending {
		return nil, errors.New("only pending orders can be cancelled")
	}

	order.Status = models.OrderCancelled
	if err := s.orderRepo.Update(order); err != nil {
		return nil, errors.New("failed to cancel order: " + err.Error())
	}

	return order, nil
}

// VoidItem voids a specific order detail item (manager action)
func (s *orderService) VoidItem(orderID string, orderDetailID string, reason string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}

	if order.Status == models.OrderPaid || order.Status == models.OrderRefunded {
		return nil, errors.New("cannot void items on paid or refunded orders")
	}

	// FIX: Also prevent voiding on cancelled orders
	if order.Status == models.OrderCancelled {
		return nil, errors.New("cannot void items on cancelled orders")
	}

	detailUUID, err := uuid.Parse(orderDetailID)
	if err != nil {
		return nil, errors.New("invalid order detail ID")
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

	found := false
	for i := range order.OrderDetails {
		if order.OrderDetails[i].ID == detailUUID {
			order.OrderDetails[i].IsVoided = true
			order.OrderDetails[i].VoidReason = reason
			found = true

			// Restore stock for non-voided quantity
			var menu models.Menu
			if err := tx.Where("id = ?", order.OrderDetails[i].MenuID).First(&menu).Error; err == nil {
				menu.DailyStock += order.OrderDetails[i].Quantity
				tx.Save(&menu)
			}
			break
		}
	}

	if !found {
		return nil, errors.New("order detail not found")
	}

	// Recalculate totals after voiding item
	order.CalculateTotal()

	if err := tx.Save(order).Error; err != nil {
		return nil, errors.New("failed to update order after voiding item")
	}

	success = true
	tx.Commit()

	return order, nil
}

// VoidOrder voids an entire order (manager action)
func (s *orderService) VoidOrder(orderID string, reason string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}

	if order.Status == models.OrderPaid || order.Status == models.OrderRefunded {
		return nil, errors.New("cannot void paid or refunded orders")
	}

	// FIX: Also prevent voiding already cancelled orders
	if order.Status == models.OrderCancelled {
		return nil, errors.New("cannot void already cancelled orders")
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

	// Void all items and restore stock
	for i := range order.OrderDetails {
		if !order.OrderDetails[i].IsVoided {
			order.OrderDetails[i].IsVoided = true
			order.OrderDetails[i].VoidReason = reason

			// Restore stock
			var menu models.Menu
			if err := tx.Where("id = ?", order.OrderDetails[i].MenuID).First(&menu).Error; err == nil {
				menu.DailyStock += order.OrderDetails[i].Quantity
				tx.Save(&menu)
			}
		}
	}

	order.Status = models.OrderCancelled
	order.CalculateTotal()

	if err := tx.Save(order).Error; err != nil {
		return nil, errors.New("failed to void order: " + err.Error())
	}

	success = true
	tx.Commit()

	return order, nil
}
