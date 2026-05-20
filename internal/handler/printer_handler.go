package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
)

type PrinterHandler struct {
	service      service.PrinterService
	orderService service.OrderService
}

func NewPrinterHandler(service service.PrinterService) *PrinterHandler {
	return &PrinterHandler{service: service}
}

func (h *PrinterHandler) SetOrderService(os service.OrderService) {
	h.orderService = os
}

func (h *PrinterHandler) PrintReceipt(c *gin.Context) {
	order, err := h.orderService.GetOrder(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", nil))
		return
	}
	receipt, err := h.service.GenerateReceipt(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate receipt", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Receipt generated", gin.H{"receipt": receipt}, nil))
}

func (h *PrinterHandler) PrintKitchen(c *gin.Context) {
	order, err := h.orderService.GetOrder(c.Param("order_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Order not found", "NOT_FOUND", nil))
		return
	}
	ticket, err := h.service.GenerateKitchenTicket(order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate kitchen ticket", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Kitchen ticket generated", gin.H{"ticket": ticket}, nil))
}
