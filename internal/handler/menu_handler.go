package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	service service.MenuService
}

func NewMenuHandler(service service.MenuService) *MenuHandler {
	return &MenuHandler{service: service}
}

func (h *MenuHandler) GetActive(c *gin.Context) {
	menus, err := h.service.GetActiveMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch menus", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Menus retrieved", menus, nil))
}

func (h *MenuHandler) GetAll(c *gin.Context) {
	menus, err := h.service.GetAllMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch all menus", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("All menus retrieved", menus, nil))
}

func (h *MenuHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "all" {
		h.GetAll(c)
		return
	}
	menu, err := h.service.GetMenuByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Menu not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Menu retrieved", menu, nil))
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req dto.CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	menu := models.Menu{
		Name: req.Name, Price: req.Price,
		DailyStock: req.DailyStock, IsActive: req.IsActive,
		IsRecipeBased: req.IsRecipeBased, Station: string(req.Station),
	}
	created, err := h.service.CreateMenu(menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create menu", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Menu created", created, nil))
}

func (h *MenuHandler) Update(c *gin.Context) {
	var req dto.UpdateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}

	// Fetch existing menu to apply partial updates cleanly
	menu, err := h.service.GetMenuByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Menu not found", "NOT_FOUND", err.Error()))
		return
	}

	if req.Name != "" {
		menu.Name = req.Name
	}
	if req.Price > 0 {
		menu.Price = req.Price
	}
	if req.Station != "" {
		menu.Station = req.Station
	}
	if req.DailyStock != nil {
		menu.DailyStock = *req.DailyStock
	}
	if req.IsActive != nil {
		menu.IsActive = *req.IsActive
	}
	if req.IsRecipeBased != nil {
		menu.IsRecipeBased = *req.IsRecipeBased
	}

	updated, err := h.service.UpdateMenu(c.Param("id"), *menu)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update menu", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Menu updated", updated, nil))
}

func (h *MenuHandler) Delete(c *gin.Context) {
	if err := h.service.DeleteMenu(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete menu", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Menu deleted", nil, nil))
}
