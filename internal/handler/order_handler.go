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

// CreateOrder godoc
// @Summary      Create a new order
// @Description  Creates a new POS order (dine_in or takeaway). Totals (subtotal, tax 11%, service charge 5%, total) are calculated server-side by Order.CalculateTotal() — do NOT include them in the request.
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateOrderRequest true "Order payload"
// @Success      201 {object} utils.Response{data=models.Order} "Order created successfully"
// @Failure      400 {object} utils.Response "Invalid request body"
// @Failure      500 {object} utils.Response "Failed to create order"
// @Router       /orders [post]
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

// GetOrders godoc
// @Summary      List all orders
// @Description  Returns all orders ordered by created_at desc. Each order includes order_details and payments (with splits).
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response{data=[]models.Order} "List of orders"
// @Failure      500 {object} utils.Response "Failed to fetch orders"
// @Router       /orders [get]
func (h *OrderHandler) GetOrders(c *gin.Context) {
	orders, err := h.orderService.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch orders", "INTERNAL_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("List of orders", orders, nil))
}

// GetOrder godoc
// @Summary      Get order by ID
// @Description  Returns a single order by UUID, including order_details and payments.splits.
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Success      200 {object} utils.Response{data=models.Order} "Order details"
// @Failure      404 {object} utils.Response "Order not found"
// @Router       /orders/{id} [get]
func (h *OrderHandler) GetOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.GetOrder(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Order details", order, nil))
}

// GetOrderByNumber godoc
// @Summary      Get order by order number
// @Description  Lookup an order using the human-readable order number (e.g. ORD-20260515-4823).
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        number path string true "Order number (e.g. ORD-20260515-4823)"
// @Success      200 {object} utils.Response{data=models.Order} "Order details"
// @Failure      404 {object} utils.Response "Order not found"
// @Router       /orders/by-number/{number} [get]
func (h *OrderHandler) GetOrderByNumber(c *gin.Context) {
	number := c.Param("number")
	order, err := h.orderService.GetOrderByNumber(number)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Order details", order, nil))
}

// UpdateOrder godoc
// @Summary      Update order notes
// @Description  Updates the free-text notes field on an existing order. Only the notes field is mutable after creation.
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Param        request body map[string]string true "Notes payload" example({"notes":"Customer requests extra napkins"})
// @Success      200 {object} utils.Response{data=models.Order} "Order updated successfully"
// @Failure      400 {object} utils.Response "Invalid request body"
// @Failure      500 {object} utils.Response "Failed to update order"
// @Router       /orders/{id} [put]
func (h *OrderHandler) UpdateOrder(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body", "BAD_REQUEST", err.Error()))
		return
	}

	order, err := h.orderService.UpdateOrderNotes(id, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update order", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Order updated successfully", order, nil))
}

// ConfirmOrder godoc
// @Summary      Confirm order and deduct stock
// @Description  Transitions order from 'pending' to 'confirmed' and deducts inventory. For recipe-based menus (is_recipe_based=true), deducts RawMaterial stock using pessimistic locking (SELECT FOR UPDATE). For stock-based menus, decrements daily_stock directly.
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Success      200 {object} utils.Response{data=models.Order} "Order confirmed"
// @Failure      500 {object} utils.Response "Failed to confirm order or insufficient stock"
// @Router       /orders/{id}/confirm [put]
func (h *OrderHandler) ConfirmOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.ConfirmOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to confirm order", "INTERNAL_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order confirmed", order, nil))
}

// CancelOrder godoc
// @Summary      Cancel a pending order (Cashier)
// @Description  Cancels an order that is still in 'pending' status. Only pending orders can be cancelled. Use VoidOrder for confirmed orders.
// @Tags         Orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Success      200 {object} utils.Response{data=models.Order} "Order cancelled"
// @Failure      400 {object} utils.Response "Order is not in pending status or other error"
// @Router       /orders/{id}/cancel [put]
func (h *OrderHandler) CancelOrder(c *gin.Context) {
	id := c.Param("id")
	order, err := h.orderService.CancelOrder(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to cancel order", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order cancelled", order, nil))
}

// VoidItem godoc
// @Summary      Void a specific line item (Manager only)
// @Description  Voids a single OrderDetail row. Sets is_voided=true, restores the menu daily_stock, and recalculates the order total via CalculateTotal(). Cannot void items on paid or refunded orders.
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Param        request body map[string]string true "Void item payload" example({"order_detail_id":"uuid","reason":"Customer changed mind"})
// @Success      200 {object} utils.Response{data=models.Order} "Item voided"
// @Failure      400 {object} utils.Response "Invalid payload or cannot void"
// @Router       /orders/{id}/void-item [post]
func (h *OrderHandler) VoidItem(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		OrderDetailID string `json:"order_detail_id" binding:"required"`
		Reason        string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body", "BAD_REQUEST", err.Error()))
		return
	}

	order, err := h.orderService.VoidItem(id, req.OrderDetailID, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to void item", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Item voided", order, nil))
}

// VoidOrder godoc
// @Summary      Void an entire order (Manager only)
// @Description  Voids all line items in the order, restores stock for each item, sets order status to 'cancelled', and recalculates totals. Cannot void paid or refunded orders.
// @Tags         Orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Param        request body map[string]string true "Void reason payload" example({"reason":"Customer walked away"})
// @Success      200 {object} utils.Response{data=models.Order} "Order voided"
// @Failure      400 {object} utils.Response "Cannot void paid/refunded order or other error"
// @Router       /orders/{id}/void [post]
func (h *OrderHandler) VoidOrder(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body", "BAD_REQUEST", err.Error()))
		return
	}

	order, err := h.orderService.VoidOrder(id, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to void order", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Order voided", order, nil))
}
