package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CashMovementHandler struct {
	service service.ShiftService
}

func NewCashMovementHandler(service service.ShiftService) *CashMovementHandler {
	return &CashMovementHandler{service: service}
}

func (h *CashMovementHandler) Add(c *gin.Context) {
	var req dto.CashMovementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	shiftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid shift ID", "VALIDATION_ERROR", nil))
		return
	}
	recordedBy, _ := c.Get("userID")
	recordedByUUID, _ := uuid.Parse(recordedBy.(string))
	cmType := models.CashMovementIn
	if req.Type == "out" {
		cmType = models.CashMovementOut
	}
	cm, err := h.service.AddCashMovement(shiftID, recordedByUUID, cmType, req.Amount, req.Reason)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to add cash movement", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Cash movement added", cm, nil))
}

func (h *CashMovementHandler) GetByShift(c *gin.Context) {
	shiftID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid shift ID", "VALIDATION_ERROR", nil))
		return
	}
	movements, err := h.service.GetCashMovements(shiftID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch movements", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Cash movements retrieved", movements, nil))
}
