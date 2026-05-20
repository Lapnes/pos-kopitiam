package handler

import (
	"net/http"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaxHandler struct {
	service service.TaxService
}

func NewTaxHandler(service service.TaxService) *TaxHandler {
	return &TaxHandler{service: service}
}

func (h *TaxHandler) GetAll(c *gin.Context) {
	// Not implemented in repo, return empty for now
	c.JSON(http.StatusOK, utils.SuccessResponse("Tax configs retrieved", []models.TaxConfig{}, nil))
}

func (h *TaxHandler) GetByID(c *gin.Context) {
	tax, err := h.service.GetTaxConfigByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Tax config not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Tax config retrieved", tax, nil))
}

func (h *TaxHandler) GetActive(c *gin.Context) {
	var branchID *uuid.UUID
	if bid := c.Query("branch_id"); bid != "" {
		id, err := uuid.Parse(bid)
		if err == nil {
			branchID = &id
		}
	}
	taxes, err := h.service.GetActiveTaxes(branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch taxes", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Active taxes retrieved", taxes, nil))
}

func (h *TaxHandler) Create(c *gin.Context) {
	var req dto.TaxConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	var branchID *uuid.UUID
	if req.BranchID != "" {
		id, _ := uuid.Parse(req.BranchID)
		branchID = &id
	}
	effectiveFrom, _ := time.Parse("2006-01-02", req.EffectiveFrom)
	tax, err := h.service.CreateTaxConfig(branchID, req.Name, models.TaxType(req.TaxType), req.Percentage, effectiveFrom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create tax config", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Tax config created", tax, nil))
}

func (h *TaxHandler) Update(c *gin.Context) {
	var req dto.TaxConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	var branchID *uuid.UUID
	if req.BranchID != "" {
		id, _ := uuid.Parse(req.BranchID)
		branchID = &id
	}
	effectiveFrom, _ := time.Parse("2006-01-02", req.EffectiveFrom)
	tax, err := h.service.UpdateTaxConfig(c.Param("id"), branchID, req.Name, models.TaxType(req.TaxType), req.Percentage, req.IsActive, effectiveFrom)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update tax config", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Tax config updated", tax, nil))
}

func (h *TaxHandler) Delete(c *gin.Context) {
	if err := h.service.DeleteTaxConfig(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete tax config", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Tax config deleted", nil, nil))
}
