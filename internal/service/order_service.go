package service

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
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
	UpdateStatus(orderID string, status string) (*models.Order, error)
	CalculateOrderTotals(order *models.Order) (subtotal, taxAmount, serviceCharge, total float64, err error)
}

type orderService struct {
	db            *gorm.DB
	orderRepo     repository.OrderRepository
	inventoryRepo repository.InventoryRepository
	shiftRepo     repository.ShiftRepository
	taxRepo       repository.TaxRepository
}

func NewOrderService(db *gorm.DB, orderRepo repository.OrderRepository, inventoryRepo repository.InventoryRepository, shiftRepo repository.ShiftRepository, taxRepo repository.TaxRepository) OrderService {
	return &orderService{
		db: db, orderRepo: orderRepo, inventoryRepo: inventoryRepo,
		shiftRepo: shiftRepo, taxRepo: taxRepo,
	}
}

func generateOrderNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("ORD-%s-%04d", time.Now().Format("20060102"), n)
}

func (s *orderService) CalculateOrderTotals(order *models.Order) (subtotal, taxAmount, serviceCharge, total float64, err error) {
	if order.OrderDetails == nil {
		return 0, 0, 0, 0, nil
	}
	for _, item := range order.OrderDetails {
		if !item.IsVoided {
			subtotal += item.Subtotal
		}
	}

	var taxes []models.TaxConfig
	if s.taxRepo != nil {
		taxes, _ = s.taxRepo.GetActiveTaxes(&order.BranchID)
	}

	taxRate := 0.11
	serviceRate := 0.05

	for _, tax := range taxes {
		if tax.TaxType == models.TaxTypeTax {
			taxRate = tax.Percentage
		}
		if tax.TaxType == models.TaxTypeService {
			serviceRate = tax.Percentage
		}
	}

	taxAmount = subtotal * taxRate
	serviceCharge = subtotal * serviceRate
	total = subtotal + taxAmount + serviceCharge - order.DiscountAmount
	if total < 0 {
		total = 0
	}

	return subtotal, taxAmount, serviceCharge, total, nil
}

// CreateOrder — FIXED: Wrap in transaction + batch menu lookups to eliminate N+1
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

	var customerID *uuid.UUID
	if req.CustomerID != "" {
		id, _ := uuid.Parse(req.CustomerID)
		customerID = &id
	}

	// === BATCH LOOKUP: Collect all menu IDs first, then query once ===
	menuIDs := make([]uuid.UUID, 0, len(req.Items))
	for _, item := range req.Items {
		id, _ := uuid.Parse(item.MenuID)
		menuIDs = append(menuIDs, id)
	}

	// Single query to fetch all menus at once (eliminates N+1)
	var menus []models.Menu
	if err := s.db.Where("id IN ? AND is_active = ?", menuIDs, true).Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch menus: %w", err)
	}

	// Build lookup map: O(n) instead of O(n²)
	menuMap := make(map[uuid.UUID]models.Menu, len(menus))
	for _, m := range menus {
		menuMap[m.ID] = m
	}

	// Identify recipe-based menus for batch recipe lookup
	recipeMenuIDs := make([]uuid.UUID, 0)
	for _, item := range req.Items {
		id, _ := uuid.Parse(item.MenuID)
		if m, ok := menuMap[id]; ok && m.IsRecipeBased {
			recipeMenuIDs = append(recipeMenuIDs, id)
		}
	}

	// Batch fetch all recipes with ingredients (2 queries total, not N)
	recipeMap := make(map[uuid.UUID]*models.Recipe)
	if len(recipeMenuIDs) > 0 {
		recipes, err := s.inventoryRepo.GetRecipesByMenuIDs(recipeMenuIDs)
		if err == nil {
			for i := range recipes {
				recipeMap[recipes[i].MenuID] = &recipes[i]
			}
		}
	}

	// Batch fetch all raw materials needed for recipes
	rawMaterialIDs := make([]uuid.UUID, 0)
	for _, recipe := range recipeMap {
		for _, ing := range recipe.Ingredients {
			rawMaterialIDs = append(rawMaterialIDs, ing.RawMaterialID)
		}
	}

	rmMap := make(map[uuid.UUID]*models.RawMaterial)
	if len(rawMaterialIDs) > 0 {
		materials, err := s.inventoryRepo.GetRawMaterialsByIDs(rawMaterialIDs)
		if err == nil {
			for i := range materials {
				rmMap[materials[i].ID] = &materials[i]
			}
		}
	}

	// === Build order details using in-memory lookups (zero DB calls) ===
	var orderDetails []models.OrderDetail
	var calculatedSubtotal float64
	for _, item := range req.Items {
		itemUUID, _ := uuid.Parse(item.MenuID)
		itemSubtotal := item.Price * float64(item.Quantity)
		calculatedSubtotal += itemSubtotal

		menu, ok := menuMap[itemUUID]
		if !ok {
			return nil, fmt.Errorf("menu item %s not found or inactive", item.MenuID)
		}

		costPrice := menu.Price * 0.5
		if menu.IsRecipeBased {
			if recipe, ok := recipeMap[itemUUID]; ok {
				costPrice = s.calculateRecipeCostBatch(recipe, rmMap)
			}
		}

		orderDetails = append(orderDetails, models.OrderDetail{
			MenuID:    itemUUID,
			MenuName:  item.MenuName,
			Quantity:  item.Quantity,
			Price:     item.Price,
			CostPrice: costPrice,
			Subtotal:  itemSubtotal,
			Notes:     item.Notes,
			Station:   menu.Station,
		})
	}

	// Calculate tax and total
	_, taxAmount, serviceCharge, total, _ := s.CalculateOrderTotals(&models.Order{
		OrderDetails:   orderDetails,
		DiscountAmount: req.DiscountAmount,
		BranchID:       branchUUID,
	})

	order := &models.Order{
		OrderNumber:    orderNumber,
		BranchID:       branchUUID,
		TableID:        tableUUID,
		EmployeeID:     userUUID,
		CustomerID:     customerID,
		OrderType:      models.OrderType(req.OrderType),
		Status:         models.OrderPending,
		Subtotal:       calculatedSubtotal,
		TaxAmount:      taxAmount,
		ServiceCharge:  serviceCharge,
		DiscountAmount: req.DiscountAmount,
		Total:          total,
		Notes:          req.Notes,
		HasReturns:     false,
		OrderDetails:   orderDetails,
	}

	// === TRANSACTION BOUNDARY: Create order + details atomically ===
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

	if err := tx.Create(order).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()
	return order, nil
}

// calculateRecipeCostBatch — uses in-memory map instead of per-ingredient DB queries
func (s *orderService) calculateRecipeCostBatch(recipe *models.Recipe, rmMap map[uuid.UUID]*models.RawMaterial) float64 {
	var totalCost float64
	for _, ingredient := range recipe.Ingredients {
		if rm, ok := rmMap[ingredient.RawMaterialID]; ok {
			totalCost += rm.CostPerUnit * ingredient.Quantity
		}
	}
	return totalCost
}

// Legacy method kept for backward compatibility (single recipe lookup)
func (s *orderService) calculateRecipeCost(recipe *models.Recipe) float64 {
	var totalCost float64
	for _, ingredient := range recipe.Ingredients {
		rm, err := s.inventoryRepo.GetRawMaterial(ingredient.RawMaterialID)
		if err == nil {
			totalCost += rm.CostPerUnit * ingredient.Quantity
		}
	}
	return totalCost
}

func (s *orderService) GetOrders() ([]models.Order, error) {
	return s.orderRepo.FindAllWithDetails()
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
		return nil, errors.New("order not found")
	}
	order.Notes = notes
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	return order, nil
}

// ConfirmOrder — FIXED: Batch all lookups inside transaction to eliminate N+1
func (s *orderService) ConfirmOrder(orderID string) (*models.Order, error) {
	// Fetch order with details OUTSIDE transaction (read-only, no lock needed)
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.Status != models.OrderPending {
		return nil, errors.New("only pending orders can be confirmed")
	}

	// Pre-collect all menu IDs and recipe menu IDs for batch lookups
	menuIDs := make([]uuid.UUID, 0, len(order.OrderDetails))
	recipeMenuIDs := make([]uuid.UUID, 0)
	for _, detail := range order.OrderDetails {
		menuIDs = append(menuIDs, detail.MenuID)
	}

	// Batch fetch all menus (1 query)
	var allMenus []models.Menu
	if err := s.db.Where("id IN ?", menuIDs).Find(&allMenus).Error; err != nil {
		return nil, fmt.Errorf("failed to batch fetch menus: %w", err)
	}
	menuMap := make(map[uuid.UUID]models.Menu)
	for _, m := range allMenus {
		menuMap[m.ID] = m
		if m.IsRecipeBased {
			recipeMenuIDs = append(recipeMenuIDs, m.ID)
		}
	}

	// Batch fetch all recipes with ingredients (1 query + 1 preload)
	recipeMap := make(map[uuid.UUID]*models.Recipe)
	if len(recipeMenuIDs) > 0 {
		recipes, err := s.inventoryRepo.GetRecipesByMenuIDs(recipeMenuIDs)
		if err == nil {
			for i := range recipes {
				recipeMap[recipes[i].MenuID] = &recipes[i]
			}
		}
	}

	// Collect all raw material IDs for batch lookup
	rawMaterialIDs := make([]uuid.UUID, 0)
	for _, recipe := range recipeMap {
		for _, ing := range recipe.Ingredients {
			rawMaterialIDs = append(rawMaterialIDs, ing.RawMaterialID)
		}
	}

	// Batch fetch all raw materials (1 query)
	rmMap := make(map[uuid.UUID]*models.RawMaterial)
	if len(rawMaterialIDs) > 0 {
		materials, err := s.inventoryRepo.GetRawMaterialsByIDs(rawMaterialIDs)
		if err == nil {
			for i := range materials {
				rmMap[materials[i].ID] = &materials[i]
			}
		}
	}

	// === TRANSACTION START: Only writes happen inside tx ===
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

	// Process each detail using in-memory lookups (zero extra DB calls)
	for _, detail := range order.OrderDetails {
		menu, ok := menuMap[detail.MenuID]
		if !ok {
			return nil, fmt.Errorf("menu %s not found", detail.MenuName)
		}

		if menu.IsRecipeBased {
			recipe, ok := recipeMap[detail.MenuID]
			if !ok || recipe == nil {
				return nil, fmt.Errorf("recipe not found for %s", detail.MenuName)
			}
			if recipe.DeletedAt.Valid {
				return nil, fmt.Errorf("recipe for %s has been deleted", detail.MenuName)
			}
			for _, ingredient := range recipe.Ingredients {
				totalDeduction := ingredient.Quantity * float64(detail.Quantity)
				rm, ok := rmMap[ingredient.RawMaterialID]
				if !ok {
					return nil, errors.New("raw material not found")
				}
				if rm.CurrentStock < totalDeduction {
					return nil, fmt.Errorf("insufficient stock for %s", rm.Name)
				}
				// Lock + update inside transaction
				var lockedRM models.RawMaterial
				if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedRM, "id = ?", ingredient.RawMaterialID).Error; err != nil {
					return nil, errors.New("failed to lock raw material")
				}
				lockedRM.CurrentStock -= totalDeduction
				if err := tx.Save(&lockedRM).Error; err != nil {
					return nil, errors.New("failed to update stock")
				}
			}
		} else {
			if menu.DailyStock < detail.Quantity {
				return nil, fmt.Errorf("insufficient daily stock for %s", menu.Name)
			}
			// Lock + update inside transaction
			var lockedMenu models.Menu
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&lockedMenu, "id = ?", detail.MenuID).Error; err != nil {
				return nil, fmt.Errorf("failed to lock menu %s", menu.Name)
			}
			lockedMenu.DailyStock -= detail.Quantity
			if lockedMenu.DailyStock < 0 {
				lockedMenu.DailyStock = 0
			}
			if err := tx.Save(&lockedMenu).Error; err != nil {
				return nil, errors.New("failed to update daily stock")
			}
		}
	}

	order.Status = models.OrderConfirmed
	if s.shiftRepo != nil {
		shift, err := s.shiftRepo.GetActiveShift(order.EmployeeID)
		if err == nil && shift != nil {
			order.ShiftID = &shift.ID
		}
	}

	if err := tx.Save(order).Error; err != nil {
		return nil, errors.New("failed to update order")
	}

	success = true
	tx.Commit()
	return order, nil
}

func (s *orderService) CancelOrder(orderID string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.Status != models.OrderPending {
		return nil, errors.New("only pending orders can be cancelled")
	}
	order.Status = models.OrderCancelled
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	return order, nil
}

// VoidItem — FIXED: Batch menu lookup + transaction boundary
func (s *orderService) VoidItem(orderID string, orderDetailID string, reason string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.Status == models.OrderPaid || order.Status == models.OrderRefunded || order.Status == models.OrderCancelled {
		return nil, errors.New("cannot void items on paid, refunded, or cancelled orders")
	}

	detailUUID, err := uuid.Parse(orderDetailID)
	if err != nil {
		return nil, errors.New("invalid order detail ID")
	}

	// Find the detail to void (outside tx — read-only)
	var targetDetail *models.OrderDetail
	for i := range order.OrderDetails {
		if order.OrderDetails[i].ID == detailUUID {
			targetDetail = &order.OrderDetails[i]
			break
		}
	}
	if targetDetail == nil {
		return nil, errors.New("order detail not found")
	}

	// Pre-fetch menu for stock restoration (single query)
	var menu models.Menu
	menuErr := s.db.Where("id = ?", targetDetail.MenuID).First(&menu).Error

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

	// Update detail inside transaction
	targetDetail.IsVoided = true
	targetDetail.VoidReason = reason

	if menuErr == nil {
		menu.DailyStock += targetDetail.Quantity
		if err := tx.Save(&menu).Error; err != nil {
			return nil, errors.New("failed to restore stock")
		}
	}

	if err := tx.Save(order).Error; err != nil {
		return nil, errors.New("failed to update order")
	}

	success = true
	tx.Commit()
	return order, nil
}

// VoidOrder — FIXED: Batch menu lookup + transaction boundary
func (s *orderService) VoidOrder(orderID string, reason string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}
	if order.Status == models.OrderPaid || order.Status == models.OrderRefunded || order.Status == models.OrderCancelled {
		return nil, errors.New("cannot void paid, refunded, or cancelled orders")
	}

	// Collect all menu IDs for batch lookup
	menuIDs := make([]uuid.UUID, 0)
	for i := range order.OrderDetails {
		if !order.OrderDetails[i].IsVoided {
			menuIDs = append(menuIDs, order.OrderDetails[i].MenuID)
		}
	}

	// Batch fetch menus (1 query)
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

	for i := range order.OrderDetails {
		if !order.OrderDetails[i].IsVoided {
			order.OrderDetails[i].IsVoided = true
			order.OrderDetails[i].VoidReason = reason

			if menu, ok := menuMap[order.OrderDetails[i].MenuID]; ok {
				menu.DailyStock += order.OrderDetails[i].Quantity
				if err := tx.Save(&menu).Error; err != nil {
					return nil, errors.New("failed to restore stock")
				}
			}
		}
	}

	order.Status = models.OrderCancelled
	if err := tx.Save(order).Error; err != nil {
		return nil, errors.New("failed to void order")
	}

	success = true
	tx.Commit()
	return order, nil
}

func (s *orderService) UpdateStatus(orderID string, status string) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return nil, errors.New("order not found")
	}

	var newStatus models.OrderStatus
	switch status {
	case "pending":
		newStatus = models.OrderPending
	case "confirmed":
		newStatus = models.OrderConfirmed
	case "preparing", "cooking":
		newStatus = models.OrderCooking
	case "ready":
		newStatus = models.OrderReady
	case "completed", "served":
		newStatus = models.OrderServed
	case "paid":
		newStatus = models.OrderPaid
	case "cancelled":
		newStatus = models.OrderCancelled
	case "refunded":
		newStatus = models.OrderRefunded
	default:
		return nil, fmt.Errorf("invalid order status: %s", status)
	}

	order.Status = newStatus
	if err := s.orderRepo.Update(order); err != nil {
		return nil, err
	}
	return order, nil
}
