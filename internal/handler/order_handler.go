package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderService
}

func NewOrderHandler(service service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	orders, err := h.service.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch orders", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Orders retrieved", orders, nil))
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	order, err := h.service.GetOrder(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order retrieved", order, nil))
}

func (h *OrderHandler) GetByNumber(c *gin.Context) {
	order, err := h.service.GetOrderByNumber(c.Param("order_number"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order retrieved", order, nil))
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	userID, _ := c.Get("userID")
	branchID, _ := c.Get("branchID")
	order, err := h.service.CreateOrder(req, userID.(string), branchID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create order", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Order created", order, nil))
}

func (h *OrderHandler) UpdateNotes(c *gin.Context) {
	var req dto.UpdateOrderNotesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	order, err := h.service.UpdateOrderNotes(c.Param("id"), req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update notes", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Notes updated", order, nil))
}

func (h *OrderHandler) Confirm(c *gin.Context) {
	order, err := h.service.ConfirmOrder(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to confirm order", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order confirmed", order, nil))
}

func (h *OrderHandler) Cancel(c *gin.Context) {
	order, err := h.service.CancelOrder(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to cancel order", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order cancelled", order, nil))
}

func (h *OrderHandler) VoidOrder(c *gin.Context) {
	var req dto.VoidOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	order, err := h.service.VoidOrder(c.Param("id"), req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to void order", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order voided", order, nil))
}

func (h *OrderHandler) VoidItem(c *gin.Context) {
	var req dto.VoidItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	order, err := h.service.VoidItem(c.Param("id"), req.OrderDetailID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to void item", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Item voided", order, nil))
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	var req dto.UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	order, err := h.service.UpdateStatus(c.Param("id"), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to update status", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order status updated", order, nil))
}
