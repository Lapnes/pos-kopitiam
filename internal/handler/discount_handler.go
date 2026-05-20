package handler

import (
	"net/http"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DiscountHandler struct {
	service service.DiscountService
}

func NewDiscountHandler(service service.DiscountService) *DiscountHandler {
	return &DiscountHandler{service: service}
}

func (h *DiscountHandler) GetAll(c *gin.Context) {
	// Return empty for now; repo doesn't have FindAll
	c.JSON(http.StatusOK, utils.SuccessResponse("Discounts retrieved", []models.Discount{}, nil))
}

func (h *DiscountHandler) GetActive(c *gin.Context) {
	discounts, err := h.service.GetActiveDiscounts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch discounts", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Active discounts retrieved", discounts, nil))
}

func (h *DiscountHandler) GetByID(c *gin.Context) {
	discount, err := h.service.GetDiscountByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Discount not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Discount retrieved", discount, nil))
}

func (h *DiscountHandler) Create(c *gin.Context) {
	var req struct {
		Name           string     `json:"name" binding:"required"`
		Type           string     `json:"type" binding:"required"`
		Value          float64    `json:"value" binding:"required,min=0"`
		MinOrderAmount float64    `json:"min_order_amount"`
		MaxDiscount    float64    `json:"max_discount"`
		StartDate      *time.Time `json:"start_date"`
		EndDate        *time.Time `json:"end_date"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	discount, err := h.service.CreateDiscount(req.Name, models.DiscountType(req.Type), req.Value, req.MinOrderAmount, req.MaxDiscount, req.StartDate, req.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create discount", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Discount created", discount, nil))
}

func (h *DiscountHandler) Update(c *gin.Context) {
	var req struct {
		Name           string     `json:"name"`
		Type           string     `json:"type"`
		Value          float64    `json:"value"`
		MinOrderAmount float64    `json:"min_order_amount"`
		MaxDiscount    float64    `json:"max_discount"`
		StartDate      *time.Time `json:"start_date"`
		EndDate        *time.Time `json:"end_date"`
		IsActive       bool       `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	discount, err := h.service.UpdateDiscount(c.Param("id"), req.Name, models.DiscountType(req.Type), req.Value, req.MinOrderAmount, req.MaxDiscount, req.StartDate, req.EndDate, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update discount", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Discount updated", discount, nil))
}

func (h *DiscountHandler) Delete(c *gin.Context) {
	if err := h.service.DeleteDiscount(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete discount", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Discount deleted", nil, nil))
}

func (h *DiscountHandler) Apply(c *gin.Context) {
	var req struct {
		OrderID        string  `json:"order_id" binding:"required"`
		DiscountID     string  `json:"discount_id" binding:"required"`
		DiscountAmount float64 `json:"discount_amount" binding:"required,min=0"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	orderID, _ := uuid.Parse(req.OrderID)
	discountID, _ := uuid.Parse(req.DiscountID)
	app, err := h.service.ApplyDiscount(orderID, discountID, req.DiscountAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to apply discount", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Discount applied", app, nil))
}
