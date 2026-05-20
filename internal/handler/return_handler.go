package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReturnHandler struct {
	service service.ReturnService
}

func NewReturnHandler(service service.ReturnService) *ReturnHandler {
	return &ReturnHandler{service: service}
}

func (h *ReturnHandler) ProcessReturn(c *gin.Context) {
	var req dto.ProcessReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	orderIDStr := c.Query("order_id")
	processedByStr := c.Query("processed_by")
	if orderIDStr == "" || processedByStr == "" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("order_id and processed_by are required", "VALIDATION_ERROR", nil))
		return
	}
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid order_id", "VALIDATION_ERROR", nil))
		return
	}
	processedBy, err := uuid.Parse(processedByStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid processed_by", "VALIDATION_ERROR", nil))
		return
	}
	orderReturn, err := h.service.ProcessReturn(orderID, processedBy, req.ReturnAmount, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Return failed", "RETURN_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Return processed", orderReturn, nil))
}
