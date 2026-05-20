package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PaymentHandler struct {
	service service.PaymentService
}

func NewPaymentHandler(service service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: service}
}

func (h *PaymentHandler) Process(c *gin.Context) {
	var req dto.ProcessPaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	orderIDStr := c.Query("order_id")
	if orderIDStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("order_id is required", "VALIDATION_ERROR", nil))
		return
	}
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid order_id", "VALIDATION_ERROR", nil))
		return
	}
	payment, err := h.service.ProcessPayment(orderID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Payment failed", "PAYMENT_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Payment processed", payment, nil))
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", "VALIDATION_ERROR", nil))
		return
	}
	payment, err := h.service.GetPaymentByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Payment not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Payment retrieved", payment, nil))
}

func (h *PaymentHandler) GetByOrderID(c *gin.Context) {
	orderID, err := uuid.Parse(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid order_id", "VALIDATION_ERROR", nil))
		return
	}
	payments, err := h.service.GetPaymentsByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch payments", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Payments retrieved", payments, nil))
}
