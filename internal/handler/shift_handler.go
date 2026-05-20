package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ShiftHandler struct {
	service service.ShiftService
}

func NewShiftHandler(service service.ShiftService) *ShiftHandler {
	return &ShiftHandler{service: service}
}

func (h *ShiftHandler) Open(c *gin.Context) {
	var req dto.OpenShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	userID, _ := c.Get("userID")
	branchID, _ := c.Get("branchID")
	userUUID, _ := uuid.Parse(userID.(string))
	branchUUID, _ := uuid.Parse(branchID.(string))
	shift, err := h.service.OpenShift(userUUID, branchUUID, req.OpeningCash)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to open shift", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Shift opened", shift, nil))
}

func (h *ShiftHandler) GetCurrent(c *gin.Context) {
	userID, _ := c.Get("userID")
	userUUID, _ := uuid.Parse(userID.(string))
	shift, err := h.service.GetCurrentShift(userUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("No active shift", "NOT_FOUND", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Current shift retrieved", shift, nil))
}

func (h *ShiftHandler) Close(c *gin.Context) {
	var req dto.CloseShiftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	userID, _ := c.Get("userID")
	userUUID, _ := uuid.Parse(userID.(string))
	shift, err := h.service.CloseShift(userUUID, req.ActualClosingCash)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to close shift", "BAD_REQUEST", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Shift closed", shift, nil))
}
