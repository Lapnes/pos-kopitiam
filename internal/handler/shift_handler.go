package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShiftHandler struct {
	shiftService service.ShiftService
}

func NewShiftHandler(shiftService service.ShiftService) *ShiftHandler {
	return &ShiftHandler{shiftService: shiftService}
}

type OpenShiftRequest struct {
	OpeningCash float64 `json:"opening_cash" binding:"required,min=0"`
}

func (h *ShiftHandler) OpenShift(c *gin.Context) {
	var req OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

	// FIX: pakai "userID" dan "branchID" (sesuai auth middleware)
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User ID not found", "UNAUTHORIZED", nil))
		return
	}
	branchIDStr, exists := c.Get("branchID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Branch ID not found", "UNAUTHORIZED", nil))
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid user ID", "UNAUTHORIZED", err.Error()))
		return
	}

	branchID, err := uuid.Parse(branchIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid branch ID", "UNAUTHORIZED", err.Error()))
		return
	}

	shift, err := h.shiftService.OpenShift(userID, branchID, req.OpeningCash)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Shift opened successfully", shift, nil))
}

func (h *ShiftHandler) GetCurrentShift(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User ID not found", "UNAUTHORIZED", nil))
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid user ID", "UNAUTHORIZED", err.Error()))
		return
	}

	shift, err := h.shiftService.GetCurrentShift(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("No active shift", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Current shift retrieved", shift, nil))
}

type CloseShiftRequest struct {
	ActualClosingCash float64 `json:"actual_closing_cash" binding:"required,min=0"`
}

func (h *ShiftHandler) CloseShift(c *gin.Context) {
	var req CloseShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("User ID not found", "UNAUTHORIZED", nil))
		return
	}
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Invalid user ID", "UNAUTHORIZED", err.Error()))
		return
	}

	shift, err := h.shiftService.CloseShift(userID, req.ActualClosingCash)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Shift closed successfully", shift, nil))
}
