package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuService service.MenuService
}

func NewMenuHandler(menuService service.MenuService) *MenuHandler {
	return &MenuHandler{menuService: menuService}
}

// GetActiveMenus godoc
// @Summary      Get all active menus
// @Description  Returns a list of all menus with is_active=true. Each menu includes its M2M categories array.
// @Tags         Menus
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response{data=[]models.Menu} "List of active menus with categories"
// @Failure      500 {object} utils.Response "Internal server error"
// @Router       /menus [get]
func (h *MenuHandler) GetActiveMenus(c *gin.Context) {
	menus, err := h.menuService.GetActiveMenus()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch menus", "INTERNAL_SERVER_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Menus retrieved successfully", menus, nil))
}
