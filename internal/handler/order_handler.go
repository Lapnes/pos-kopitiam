package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/dto"
	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req dto.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body", "BAD_REQUEST", err.Error()))
		return
	}

	userID := c.GetString("userID")
	branchID := c.GetString("branchID")

	order, err := h.orderService.CreateOrder(req, userID, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create order", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Order created successfully", order, nil))
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.orderService.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch orders", "INTERNAL_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("List of orders", orders, nil))
}

func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Order details", order, nil))
}

func (h *OrderHandler) GetOrderByNumber(c *gin.Context) {
	number := c.Param("number")
	order, err := h.orderService.GetOrderByNumber(number)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Order details", order, nil))
}

func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Order updated", nil, nil))
}

func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Order confirmed", nil, nil))
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Order cancelled", nil, nil))
}

func (h *OrderHandler) VoidItem(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Item voided", nil, nil))
}

func (h *OrderHandler) VoidOrder(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Order voided", nil, nil))
}
