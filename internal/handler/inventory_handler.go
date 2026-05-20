package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	service service.InventoryService
}

func NewInventoryHandler(service service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) GetRawMaterial(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", "VALIDATION_ERROR", nil))
		return
	}
	rm, err := h.service.GetRawMaterial(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Raw material not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Raw material retrieved", rm, nil))
}

func (h *InventoryHandler) GetAllRawMaterials(c *gin.Context) {
	rms, err := h.service.GetAllRawMaterials()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch raw materials", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Raw materials retrieved", rms, nil))
}

func (h *InventoryHandler) CreateRawMaterial(c *gin.Context) {
	var rm models.RawMaterial
	if err := c.ShouldBindJSON(&rm); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	if err := h.service.CreateRawMaterial(&rm); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create raw material", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Raw material created", rm, nil))
}

func (h *InventoryHandler) UpdateRawMaterial(c *gin.Context) {
	var rm models.RawMaterial
	if err := c.ShouldBindJSON(&rm); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", "VALIDATION_ERROR", nil))
		return
	}
	rm.ID = id
	if err := h.service.UpdateRawMaterial(&rm); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update raw material", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Raw material updated", rm, nil))
}

func (h *InventoryHandler) DeleteRawMaterial(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", "VALIDATION_ERROR", nil))
		return
	}
	if err := h.service.DeleteRawMaterial(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete raw material", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Raw material deleted", nil, nil))
}

func (h *InventoryHandler) GetRecipeByMenu(c *gin.Context) {
	menuID, err := uuid.Parse(c.Param("menu_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid menu_id", "VALIDATION_ERROR", nil))
		return
	}
	recipe, err := h.service.GetRecipeByMenuID(menuID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Recipe not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Recipe retrieved", recipe, nil))
}

func (h *InventoryHandler) UpsertRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	result, err := h.service.UpsertRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to upsert recipe", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Recipe upserted", result, nil))
}

func (h *InventoryHandler) UpdateRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", "VALIDATION_ERROR", nil))
		return
	}
	recipe.ID = id
	result, err := h.service.UpdateRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update recipe", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Recipe updated", result, nil))
}

func (h *InventoryHandler) DeleteRecipe(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid ID", "VALIDATION_ERROR", nil))
		return
	}
	if err := h.service.DeleteRecipe(id); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete recipe", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Recipe deleted", nil, nil))
}

func (h *InventoryHandler) StockEstimation(c *gin.Context) {
	estimation, err := h.service.GetStockEstimation()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch stock estimation", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Stock estimation retrieved", estimation, nil))
}
