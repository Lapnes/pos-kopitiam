package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch categories", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Categories retrieved", categories, nil))
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	category, err := h.service.GetCategoryByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Category not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Category retrieved", category, nil))
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	category, err := h.service.CreateCategory(req.Name, req.SortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create category", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Category created", category, nil))
}

func (h *CategoryHandler) Update(c *gin.Context) {
	var req dto.CategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	category, err := h.service.UpdateCategory(c.Param("id"), req.Name, req.SortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update category", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Category updated", category, nil))
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	if err := h.service.DeleteCategory(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete category", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Category deleted", nil, nil))
}
