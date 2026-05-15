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

// GetReceiptPayload godoc
// @Summary      Generate receipt payload for printing
// @Description  Returns a structured JSON payload ready for a receipt printer. Includes store header, order details, totals, and the payment splits array (normalized — iterates payment.splits for method/amount since those fields were moved out of the Payment header).
// @Tags         Printers
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Success      200 {object} utils.Response{data=map[string]interface{}} "Receipt payload"
// @Failure      400 {object} utils.Response "Invalid order UUID"
// @Failure      404 {object} utils.Response "Order not found or failed to generate receipt"
// @Router       /printers/receipt/{id} [get]
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

// GetKitchenTicketPayload godoc
// @Summary      Generate kitchen ticket payload
// @Description  Returns a structured JSON payload for the Kitchen Display System (KDS) or kitchen printer. Includes only the items assigned to the kitchen station (station=kitchen), along with order notes and table/order type info.
// @Tags         Printers
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Order UUID"
// @Success      200 {object} utils.Response{data=map[string]interface{}} "Kitchen ticket payload"
// @Failure      400 {object} utils.Response "Invalid order UUID"
// @Failure      404 {object} utils.Response "Order not found or failed to generate ticket"
// @Router       /printers/kitchen/{id} [get]
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
