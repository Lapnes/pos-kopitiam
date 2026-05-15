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

// OpenShift godoc
// @Summary      Open a cashier shift
// @Description  Opens a new shift for the authenticated user. Only one shift can be open per user at a time. Records opening cash balance.
// @Tags         Shifts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body OpenShiftRequest true "Opening cash amount"
// @Success      201 {object} utils.Response{data=models.Shift} "Shift opened successfully"
// @Failure      400 {object} utils.Response "Shift already open or invalid payload"
// @Failure      401 {object} utils.Response "Unauthorized — user ID missing"
// @Router       /shifts/open [post]
func (h *ShiftHandler) OpenShift(c *gin.Context) {
	var req OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

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

// GetCurrentShift godoc
// @Summary      Get the current open shift
// @Description  Returns the currently active shift for the authenticated user. Returns 404 if no shift is open.
// @Tags         Shifts
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response{data=models.Shift} "Current shift details"
// @Failure      401 {object} utils.Response "Unauthorized"
// @Failure      404 {object} utils.Response "No active shift found"
// @Router       /shifts/current [get]
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

// CloseShift godoc
// @Summary      Close the current shift
// @Description  Closes the active shift for the authenticated user. Calculates variance between expected cash (opening + sales) and actual counted cash. Returns shift summary with total_sales, variance, and total_transactions.
// @Tags         Shifts
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CloseShiftRequest true "Actual closing cash count"
// @Success      200 {object} utils.Response{data=models.Shift} "Shift closed successfully with summary"
// @Failure      400 {object} utils.Response "No open shift or invalid payload"
// @Failure      401 {object} utils.Response "Unauthorized"
// @Router       /shifts/close [post]
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
