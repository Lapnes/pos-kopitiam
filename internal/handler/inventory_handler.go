package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type InventoryHandler struct {
	inventoryService service.InventoryService
}

func NewInventoryHandler(inventoryService service.InventoryService) *InventoryHandler {
	return &InventoryHandler{inventoryService: inventoryService}
}

// RAW MATERIALS

func (h *InventoryHandler) CreateRawMaterial(c *gin.Context) {
	var rm models.RawMaterial
	if err := c.ShouldBindJSON(&rm); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

	if err := h.inventoryService.CreateRawMaterial(&rm); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Raw material created successfully", rm, nil))
}

func (h *InventoryHandler) UpdateRawMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid raw material ID", "BAD_REQUEST", err.Error()))
		return
	}

	var rm models.RawMaterial
	if err := c.ShouldBindJSON(&rm); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}
	rm.ID = id

	if err := h.inventoryService.UpdateRawMaterial(&rm); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Raw material updated successfully", rm, nil))
}

func (h *InventoryHandler) GetRawMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid raw material ID", "BAD_REQUEST", err.Error()))
		return
	}

	rm, err := h.inventoryService.GetRawMaterial(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Raw material not found", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Raw material retrieved", rm, nil))
}

func (h *InventoryHandler) DeleteRawMaterial(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid raw material ID", "BAD_REQUEST", err.Error()))
		return
	}

	if err := h.inventoryService.DeleteRawMaterial(id); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to delete raw material", "BAD_REQUEST", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Raw material deleted successfully", nil, nil))
}

// RECIPES

func (h *InventoryHandler) CreateRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

	if err := h.inventoryService.CreateRecipe(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Recipe created successfully", recipe, nil))
}

func (h *InventoryHandler) UpdateRecipe(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid recipe ID", "BAD_REQUEST", err.Error()))
		return
	}

	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}
	recipe.ID = id

	if err := h.inventoryService.UpdateRecipe(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Recipe updated successfully", recipe, nil))
}

func (h *InventoryHandler) GetRecipe(c *gin.Context) {
	menuIDStr := c.Param("menu_id") // We usually fetch by menu_id
	menuID, err := uuid.Parse(menuIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid menu ID", "BAD_REQUEST", err.Error()))
		return
	}

	recipe, err := h.inventoryService.GetRecipeByMenuID(menuID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Recipe not found", "NOT_FOUND", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Recipe retrieved", recipe, nil))
}

func (h *InventoryHandler) DeleteRecipe(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid recipe ID", "BAD_REQUEST", err.Error()))
		return
	}

	if err := h.inventoryService.DeleteRecipe(id); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Failed to delete recipe", "BAD_REQUEST", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Recipe deleted successfully", nil, nil))
}
