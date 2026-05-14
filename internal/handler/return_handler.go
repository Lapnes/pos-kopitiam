package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReturnHandler struct {
	returnService service.ReturnService
}

func NewReturnHandler(returnService service.ReturnService) *ReturnHandler {
	return &ReturnHandler{returnService: returnService}
}

type ProcessReturnRequest struct {
	ReturnAmount float64 `json:"return_amount" binding:"required,min=0"`
	Reason       string  `json:"reason" binding:"required"`
}

func (h *ReturnHandler) ProcessReturn(c *gin.Context) {
	var req ProcessReturnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

	orderIDStr := c.Param("id")
	orderID, err := uuid.Parse(orderIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid order ID", "BAD_REQUEST", err.Error()))
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User ID not found in context", "UNAUTHORIZED", nil))
		return
	}
	processedBy, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid user ID", "UNAUTHORIZED", err.Error()))
		return
	}

	orderReturn, err := h.returnService.ProcessReturn(orderID, processedBy, req.ReturnAmount, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Return processed successfully", orderReturn, nil))
}
