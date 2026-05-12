package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/dto"
	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type OrderService interface {
	CreateOrder(req dto.CreateOrderRequest, userID string, branchID string) (*models.Order, error)
	GetOrders() ([]models.Order, error)
	GetOrder(id string) (*models.Order, error)
	GetOrderByNumber(orderNumber string) (*models.Order, error)
	ConfirmOrder(orderID string) (*models.Order, error)
}

type orderService struct {
	db            *gorm.DB
	orderRepo     repository.OrderRepository
	inventoryRepo repository.InventoryRepository
	redisClient   *redis.Client
}

func NewOrderService(db *gorm.DB, orderRepo repository.OrderRepository, inventoryRepo repository.InventoryRepository, redisClient *redis.Client) OrderService {
	return &orderService{
		db:            db,
		orderRepo:     orderRepo,
		inventoryRepo: inventoryRepo,
		redisClient:   redisClient,
	}
}

func (s *orderService) CreateOrder(req dto.CreateOrderRequest, userID string, branchID string) (*models.Order, error) {
	// 1. Validate items and stock (Simplified for now)
	if len(req.Items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	// Generate Order Number ORD-YYYYMMDD-XXXX
	orderNumber := fmt.Sprintf("ORD-%s-%d", time.Now().Format("20060102"), time.Now().Unix()%10000)

	branchUUID, _ := uuid.Parse(branchID)
	userUUID, _ := uuid.Parse(userID)
	var tableUUID *uuid.UUID
	if req.TableID != "" {
		id, _ := uuid.Parse(req.TableID)
		tableUUID = &id
	}

	var subtotal float64
	var orderDetails []models.OrderDetail
	for _, item := range req.Items {
		itemUUID, _ := uuid.Parse(item.MenuID)
		itemSubtotal := item.Price * float64(item.Quantity)
		subtotal += itemSubtotal
		
		orderDetails = append(orderDetails, models.OrderDetail{
			MenuID:   itemUUID,
			MenuName: item.MenuName,
			Quantity: item.Quantity,
			Price:    item.Price,
			Subtotal: itemSubtotal,
			Notes:    item.Notes,
			Station:  models.StationKitchen, // Should be fetched from Menu repo
		})
	}

	taxAmount := subtotal * 0.11 // 11% PPN
	serviceCharge := subtotal * 0.05 // 5% Service Charge
	total := subtotal + taxAmount + serviceCharge - req.DiscountAmount

	order := &models.Order{
		OrderNumber:    orderNumber,
		BranchID:       branchUUID,
		TableID:        tableUUID,
		EmployeeID:     userUUID,
		CustomerName:   req.CustomerName,
		CustomerPhone:  req.CustomerPhone,
		OrderType:      models.OrderType(req.OrderType),
		Status:         models.OrderPending,
		Subtotal:       subtotal,
		TaxAmount:      taxAmount,
		ServiceCharge:  serviceCharge,
		DiscountAmount: req.DiscountAmount,
		Total:          total,
		Notes:          req.Notes,
		OrderDetails:   orderDetails,
	}

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

	// Flag to ensure we don't commit if an error occurs
	success := false
	defer func() {
		if !success {
			tx.Rollback()
		}
	}()

	for _, detail := range order.OrderDetails {
		// Get Recipe
		var recipe models.Recipe
		err := tx.Preload("Ingredients").Where("menu_id = ?", detail.MenuID).First(&recipe).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("failed to fetch recipe for menu %s: %w", detail.MenuName, err)
		}

		if err == nil {
			// Deduct Raw Materials based on BOM
			for _, ingredient := range recipe.Ingredients {
				totalDeduction := ingredient.Quantity * float64(detail.Quantity)

				var rawMaterial models.RawMaterial
				if err := tx.Where("id = ?", ingredient.RawMaterialID).First(&rawMaterial).Error; err != nil {
					return nil, errors.New("raw material not found")
				}

				if rawMaterial.CurrentStock < totalDeduction {
					return nil, fmt.Errorf("insufficient stock for raw material: %s", rawMaterial.Name)
				}

				rawMaterial.CurrentStock -= totalDeduction
				if err := tx.Save(&rawMaterial).Error; err != nil {
					return nil, errors.New("failed to update raw material stock")
				}
			}
		}

		// Decrement Menu Daily Stock (Snapshot)
		var menu models.Menu
		if err := tx.Where("id = ?", detail.MenuID).First(&menu).Error; err == nil {
			menu.DailyStock -= detail.Quantity
			if menu.DailyStock < 0 {
				menu.DailyStock = 0
			}
			tx.Save(&menu)
		}
	}

	// Update order status
	order.Status = models.OrderConfirmed
	if err := tx.Save(order).Error; err != nil {
		return nil, errors.New("failed to update order status")
	}

	success = true
	tx.Commit()

	return order, nil
}
