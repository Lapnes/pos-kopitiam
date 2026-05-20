package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CatalogHandler struct {
	db               *gorm.DB
	menuService      service.MenuService
	analyticsService service.AnalyticsService
	orderService     service.OrderService
	paymentService   service.PaymentService
}

func NewCatalogHandler(
	db *gorm.DB,
	menuService service.MenuService,
	analyticsService service.AnalyticsService,
	orderService service.OrderService,
	paymentService service.PaymentService,
) *CatalogHandler {
	return &CatalogHandler{
		db:               db,
		menuService:      menuService,
		analyticsService: analyticsService,
		orderService:     orderService,
		paymentService:   paymentService,
	}
}

type CatalogMenuResponse struct {
	models.Menu
	IsBestSeller bool `json:"is_best_seller"`
}

func (h *CatalogHandler) GetMenus(c *gin.Context) {
	menus, err := h.menuService.GetActiveMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch menus", "SERVER_ERROR", err.Error()))
		return
	}

	// Fetch best sellers from analytics for the past 30 days
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -30)
	bestSellers, err := h.analyticsService.GetBestSellers(10, startDate, endDate)
	if err != nil {
		// Log error but proceed without marking best sellers
		bestSellers = nil
	}

	bestSellerMap := make(map[string]bool)
	if bestSellers != nil {
		for _, bs := range bestSellers {
			if name, ok := bs["menu_name"].(string); ok {
				bestSellerMap[name] = true
			}
		}
	}

	response := make([]CatalogMenuResponse, 0, len(menus))
	for _, m := range menus {
		_, exists := bestSellerMap[m.Name]
		response = append(response, CatalogMenuResponse{
			Menu:         m,
			IsBestSeller: exists,
		})
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Catalog menus retrieved", response, nil))
}

type CatalogCheckoutRequest struct {
	Items []struct {
		MenuID   string `json:"menu_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=1"`
		Notes    string `json:"notes"`
	} `json:"items" binding:"required,min=1"`
	PaymentMethod string `json:"payment_method" binding:"required,oneof=online qris bank_transfer gopay debit_card"`
	TableNumber   string `json:"table_number" binding:"required"`
}

func getDeterministicTableUUID(tableNum string) string {
	if tableNum == "" {
		return ""
	}
	return uuid.NewSHA1(uuid.NameSpaceDNS, []byte("table-"+tableNum)).String()
}

func (h *CatalogHandler) Checkout(c *gin.Context) {
	var req CatalogCheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body", "VALIDATION_ERROR", err.Error()))
		return
	}

	// Retrieve branch and user details from DB
	var branch models.Branch
	if err := h.db.First(&branch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Database error: branch not found", "SERVER_ERROR", err.Error()))
		return
	}

	var cashier models.User
	if err := h.db.Where("role = ?", models.RoleCashier).First(&cashier).Error; err != nil {
		// Fallback to first user
		if err := h.db.First(&cashier).Error; err != nil {
			c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Database error: user not found", "SERVER_ERROR", err.Error()))
			return
		}
	}

	// Create order items by resolving menu details from database
	orderItems := make([]dto.OrderItemRequest, 0, len(req.Items))
	for _, item := range req.Items {
		menu, err := h.menuService.GetMenuByID(item.MenuID)
		if err != nil {
			c.JSON(http.StatusBadRequest, utils.ErrorResponse(fmt.Sprintf("Invalid menu_id: %s", item.MenuID), "VALIDATION_ERROR", err.Error()))
			return
		}
		orderItems = append(orderItems, dto.OrderItemRequest{
			MenuID:   menu.ID.String(),
			MenuName: menu.Name,
			Quantity: item.Quantity,
			Price:    menu.Price,
			Notes:    item.Notes,
		})
	}

	// Create the order request DTO
	orderReq := dto.CreateOrderRequest{
		OrderType:      "dine_in",
		TableID:        getDeterministicTableUUID(req.TableNumber),
		Items:          orderItems,
		Notes:          fmt.Sprintf("Meja %s | Self-ordered from catalog tablet", req.TableNumber),
	}

	// 1. Create order (pending status)
	order, err := h.orderService.CreateOrder(orderReq, cashier.ID.String(), branch.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create order", "SERVER_ERROR", err.Error()))
		return
	}

	// 2. Confirm order stock
	confirmedOrder, err := h.orderService.ConfirmOrder(order.ID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to confirm order stock", "SERVER_ERROR", err.Error()))
		return
	}

	// Map UI payment methods to GORM enums
	var dbPaymentMethod string
	if req.PaymentMethod == "qris" {
		dbPaymentMethod = "qris"
	} else if req.PaymentMethod == "bank_transfer" {
		dbPaymentMethod = "transfer"
	} else if req.PaymentMethod == "gopay" {
		dbPaymentMethod = "e_wallet"
	} else { // online / debit_card
		dbPaymentMethod = "debit_card"
	}

	// 3. Process payment
	paymentReq := dto.ProcessPaymentRequest{
		ChangeAmount: 0,
		Splits: []dto.PaymentSplitRequest{
			{
				PaymentMethod: dbPaymentMethod,
				Amount:        confirmedOrder.Total,
			},
		},
	}

	payment, err := h.paymentService.ProcessPayment(confirmedOrder.ID, paymentReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to process payment", "SERVER_ERROR", err.Error()))
		return
	}

	// Log catalog checkout to audit logs
	go func() {
		log := models.AuditLog{
			UserID:    cashier.ID,
			Action:    string(models.AuditCreate),
			Entity:    "orders",
			EntityID:  confirmedOrder.ID,
			OldValue:  "null",
			NewValue:  fmt.Sprintf(`{"order_id":"%s","total":%f,"table_number":"%s"}`, confirmedOrder.ID.String(), confirmedOrder.Total, req.TableNumber),
			Timestamp: time.Now(),
		}
		log.ID = uuid.New()
		log.CreatedAt = time.Now()
		log.UpdatedAt = time.Now()
		h.db.Create(&log)
	}()

	c.JSON(http.StatusOK, utils.SuccessResponse("Checkout successful", gin.H{
		"order_id":   confirmedOrder.ID.String(),
		"total":      confirmedOrder.Total,
		"snap_token": payment.SnapToken,
		"snap_url":   payment.SnapURL,
	}, nil))
}
