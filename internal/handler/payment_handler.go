package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// PaymentHandler handles payment processing for orders.
// It takes db directly since payment logic is simple enough not to warrant a full service layer yet.
type PaymentHandler struct {
	db *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

// PaymentSplitRequest represents one payment method in a split payment.
type PaymentSplitRequest struct {
	PaymentMethod   string  `json:"payment_method" binding:"required"`
	Amount          float64 `json:"amount" binding:"required,min=0"`
	ReferenceNumber string  `json:"reference_number"`
}

// ProcessPaymentRequest is the normalized payment payload.
// top-level fields hold totals; per-method detail lives in splits[].
type ProcessPaymentRequest struct {
	TotalPaid    float64               `json:"total_paid" binding:"required,min=0"`
	ChangeAmount float64               `json:"change_amount"`
	Splits       []PaymentSplitRequest `json:"splits" binding:"required,min=1"`
}

// ProcessPayment godoc
// @Summary      Process payment for an order
// @Description  Creates a normalized Payment record with one or more PaymentSplit rows in a single transaction. Supports single-method (cash, QRIS) and split payments (e.g. part cash, part QRIS). The Payment header stores totals; per-method details live in splits[]. Valid payment_method values: cash | qris | debit_card | credit_card | transfer | e_wallet.
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Param        request body ProcessPaymentRequest true "Payment payload with splits"
// @Success      201 {object} utils.Response{data=models.Payment} "Payment processed successfully"
// @Failure      400 {object} utils.Response "Invalid payload or order not found"
// @Failure      500 {object} utils.Response "Transaction failed"
// @Router       /orders/{id}/payments [post]
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid order ID", "BAD_REQUEST", err.Error()))
		return
	}

	var req ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid payment request", "BAD_REQUEST", err.Error()))
		return
	}

	// Verify the order exists
	var order models.Order
	if err := h.db.First(&order, "id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", err.Error()))
		return
	}

	// Build Payment + PaymentSplit models
	payment := models.Payment{
		OrderID:      orderID,
		TotalPaid:    req.TotalPaid,
		ChangeAmount: req.ChangeAmount,
		Status:       "success",
	}

	splits := make([]models.PaymentSplit, 0, len(req.Splits))
	for _, s := range req.Splits {
		splits = append(splits, models.PaymentSplit{
			PaymentMethod:   models.PaymentMethod(s.PaymentMethod),
			Amount:          s.Amount,
			ReferenceNumber: s.ReferenceNumber,
		})
	}

	// Persist in a transaction: create Payment then its splits
	tx := h.db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to start transaction", "INTERNAL_ERROR", tx.Error.Error()))
		return
	}

	if err := tx.Create(&payment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create payment", "INTERNAL_ERROR", err.Error()))
		return
	}

	// Attach PaymentID to each split now that payment.ID is populated
	for i := range splits {
		splits[i].PaymentID = payment.ID
	}

	if err := tx.Create(&splits).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create payment splits", "INTERNAL_ERROR", err.Error()))
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to commit payment", "INTERNAL_ERROR", err.Error()))
		return
	}

	// Attach splits to response
	payment.Splits = splits

	c.JSON(http.StatusCreated, utils.SuccessResponse("Payment processed successfully", payment, nil))
}
