package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PrinterHandler struct {
	printerService service.PrinterService
}

func NewPrinterHandler(printerService service.PrinterService) *PrinterHandler {
	return &PrinterHandler{printerService: printerService}
}

func (h *PrinterHandler) GetReceiptPayload(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid order ID", "BAD_REQUEST", err.Error()))
		return
	}

	payload, err := h.printerService.GenerateReceipt(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to generate receipt", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Receipt payload generated", payload, nil))
}

func (h *PrinterHandler) GetKitchenTicketPayload(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid order ID", "BAD_REQUEST", err.Error()))
		return
	}

	payload, err := h.printerService.GenerateKitchenTicket(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Failed to generate kitchen ticket", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Kitchen ticket payload generated", payload, nil))
}
