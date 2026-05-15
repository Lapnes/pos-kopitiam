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

// ── RAW MATERIALS ─────────────────────────────────────────────────────────────

// CreateRawMaterial godoc
// @Summary      Create a raw material
// @Description  Creates a new raw material entry (e.g. Arabica Coffee Beans, 5000g). Used as BOM ingredients in Recipes.
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.RawMaterial true "Raw material payload"
// @Success      201 {object} utils.Response{data=models.RawMaterial} "Raw material created successfully"
// @Failure      400 {object} utils.Response "Invalid payload or validation error"
// @Router       /master/raw-materials [post]
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

// UpdateRawMaterial godoc
// @Summary      Update a raw material
// @Description  Updates name, unit, stock levels, or cost of an existing raw material by UUID.
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Raw Material UUID"
// @Param        request body models.RawMaterial true "Updated raw material payload"
// @Success      200 {object} utils.Response{data=models.RawMaterial} "Raw material updated successfully"
// @Failure      400 {object} utils.Response "Invalid UUID or payload"
// @Router       /master/raw-materials/{id} [put]
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

// GetRawMaterial godoc
// @Summary      Get a raw material by ID
// @Description  Fetches a single raw material record including current_stock and minimum_stock.
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Raw Material UUID"
// @Success      200 {object} utils.Response{data=models.RawMaterial} "Raw material details"
// @Failure      400 {object} utils.Response "Invalid UUID"
// @Failure      404 {object} utils.Response "Raw material not found"
// @Router       /master/raw-materials/{id} [get]
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

// DeleteRawMaterial godoc
// @Summary      Delete a raw material
// @Description  Soft-deletes a raw material by UUID. Ensure no active recipes reference this material before deleting.
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Raw Material UUID"
// @Success      200 {object} utils.Response "Raw material deleted successfully"
// @Failure      400 {object} utils.Response "Invalid UUID or deletion failed"
// @Router       /master/raw-materials/{id} [delete]
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

// ── RECIPES (BOM) ─────────────────────────────────────────────────────────────

// CreateRecipe godoc
// @Summary      Create a recipe (BOM) for a menu item
// @Description  Links a set of raw material ingredients to a menu (Bill of Materials). If a recipe already exists for the given menu_id, it performs an upsert (update). The menu must have is_recipe_based=true for ConfirmOrder to use this recipe for stock deduction.
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body models.Recipe true "Recipe payload with ingredients array"
// @Success      201 {object} utils.Response{data=models.Recipe} "Recipe created successfully"
// @Failure      400 {object} utils.Response "Invalid payload or empty ingredients"
// @Router       /master/recipes [post]
func (h *InventoryHandler) CreateRecipe(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

	// FIX: Use UpsertRecipe instead of CreateRecipe to handle existing recipes
	createdRecipe, err := h.inventoryService.UpsertRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Recipe created successfully", createdRecipe, nil))
}

// UpdateRecipe godoc
// @Summary      Update a recipe
// @Description  Replaces the instructions and ingredients of an existing recipe by Recipe UUID.
// @Tags         Inventory
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Recipe UUID"
// @Param        request body models.Recipe true "Updated recipe payload"
// @Success      200 {object} utils.Response{data=models.Recipe} "Recipe updated successfully"
// @Failure      400 {object} utils.Response "Invalid UUID or payload"
// @Router       /master/recipes/{id} [put]
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

	updatedRecipe, err := h.inventoryService.UpdateRecipe(&recipe)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse(err.Error(), "BAD_REQUEST", nil))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Recipe updated successfully", updatedRecipe, nil))
}

// GetRecipe godoc
// @Summary      Get recipe by Menu ID
// @Description  Fetches a recipe and its ingredient list using the associated menu's UUID (not the recipe ID).
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        menu_id path string true "Menu UUID"
// @Success      200 {object} utils.Response{data=models.Recipe} "Recipe with ingredients"
// @Failure      400 {object} utils.Response "Invalid menu UUID"
// @Failure      404 {object} utils.Response "Recipe not found for this menu"
// @Router       /master/recipes/{menu_id} [get]
func (h *InventoryHandler) GetRecipe(c *gin.Context) {
	menuIDStr := c.Param("menu_id")
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

// DeleteRecipe godoc
// @Summary      Delete a recipe
// @Description  Soft-deletes a recipe and its ingredient rows by Recipe UUID.
// @Tags         Inventory
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Recipe UUID"
// @Success      200 {object} utils.Response "Recipe deleted successfully"
// @Failure      400 {object} utils.Response "Invalid UUID or deletion failed"
// @Router       /master/recipes/{id} [delete]
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
