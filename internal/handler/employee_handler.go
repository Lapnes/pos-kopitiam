package handler

import (
	"net/http"
	"strings"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmployeeHandler struct {
	service service.EmployeeService
}

func NewEmployeeHandler(service service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) GetAll(c *gin.Context) {
	employees, err := h.service.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch employees", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Employees retrieved", employees, nil))
}

func (h *EmployeeHandler) GetByID(c *gin.Context) {
	employee, err := h.service.GetEmployeeByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Employee not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Employee retrieved", employee, nil))
}

func (h *EmployeeHandler) Create(c *gin.Context) {
	var req dto.CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request structure", "VALIDATION_ERROR", err.Error()))
		return
	}

	req.Role = strings.ToLower(strings.TrimSpace(req.Role))

	// Semantic validation for Create
	if req.Role != "manager" && req.Role != "cashier" && req.Role != "kitchen" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid role. Role must be 'manager', 'cashier', or 'kitchen'", "VALIDATION_ERROR", nil))
		return
	}
	if len(req.PINCode) != 6 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("PIN code must be exactly 6 digits", "VALIDATION_ERROR", nil))
		return
	}
	if req.Password != "" && len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Password must be at least 6 characters long", "VALIDATION_ERROR", nil))
		return
	}

	branchUUID, _ := uuid.Parse(req.BranchID)

	// FIX: Pass password from DTO to service (was passing empty string)
	employee, err := h.service.CreateEmployee(
		req.Name, req.Email, req.Password, models.Role(req.Role), req.PINCode, branchUUID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create employee", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Employee created", employee, nil))
}

func (h *EmployeeHandler) Update(c *gin.Context) {
	var req dto.UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request structure", "VALIDATION_ERROR", err.Error()))
		return
	}

	if req.Role != "" {
		req.Role = strings.ToLower(strings.TrimSpace(req.Role))
	}

	// Semantic validation for Update
	if req.Role != "" && req.Role != "manager" && req.Role != "cashier" && req.Role != "kitchen" {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid role. Role must be 'manager', 'cashier', or 'kitchen'", "VALIDATION_ERROR", nil))
		return
	}
	if req.PINCode != "" && len(req.PINCode) != 6 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("PIN code must be exactly 6 digits", "VALIDATION_ERROR", nil))
		return
	}
	if req.Password != "" && len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Password must be at least 6 characters long", "VALIDATION_ERROR", nil))
		return
	}

	var isActive bool
	if req.IsActive != nil {
		isActive = *req.IsActive
	}
	// FIX: Pass password and pin from DTO to service for updates
	employee, err := h.service.UpdateEmployee(
		c.Param("id"), req.Name, req.Email, req.Password, models.Role(req.Role), req.PINCode, isActive,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update employee", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Employee updated", employee, nil))
}

func (h *EmployeeHandler) Delete(c *gin.Context) {
	if err := h.service.DeleteEmployee(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete employee", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Employee deleted", nil, nil))
}
