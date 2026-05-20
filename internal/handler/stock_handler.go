package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StockHandler struct {
	service service.StockService
}

func NewStockHandler(service service.StockService) *StockHandler {
	return &StockHandler{service: service}
}

func (h *StockHandler) Adjust(c *gin.Context) {
	var req dto.StockAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	rawMaterialID, err := uuid.Parse(req.RawMaterialID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid raw_material_id", "VALIDATION_ERROR", nil))
		return
	}
	adjustedBy, _ := c.Get("userID")
	adjustedByUUID, _ := uuid.Parse(adjustedBy.(string))
	adjustment, err := h.service.AdjustStock(rawMaterialID, req.QuantityAfter, req.Reason, adjustedByUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to adjust stock", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Stock adjusted", adjustment, nil))
}

func (h *StockHandler) GetByRawMaterial(c *gin.Context) {
	rawMaterialID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", "VALIDATION_ERROR", nil))
		return
	}
	adjustments, err := h.service.GetAdjustmentsByRawMaterial(rawMaterialID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch adjustments", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Adjustments retrieved", adjustments, nil))
}
