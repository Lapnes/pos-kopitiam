package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmployeeHandler struct {
	employeeService service.EmployeeService
}

func NewEmployeeHandler(employeeService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{employeeService: employeeService}
}

type CreateEmployeeRequest struct {
	BranchID string      `json:"branch_id" binding:"required"`
	Name     string      `json:"name" binding:"required"`
	Email    string      `json:"email"`
	Password string      `json:"password" binding:"required"`
	PINCode  string      `json:"pin_code" binding:"required"`
	Role     models.Role `json:"role" binding:"required"`
	IsActive bool        `json:"is_active"`
}

type UpdateEmployeeRequest struct {
	Name     string      `json:"name"`
	Email    string      `json:"email"`
	Password string      `json:"password"`
	PINCode  string      `json:"pin_code"`
	Role     models.Role `json:"role"`
	IsActive *bool       `json:"is_active"`
}

func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	employees, err := h.employeeService.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch employees", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Employees retrieved successfully", employees, nil))
}

func (h *EmployeeHandler) CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

	branchUUID, err := uuid.Parse(req.BranchID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid branch ID", "BAD_REQUEST", err.Error()))
		return
	}

	employee, err := h.employeeService.CreateEmployee(branchUUID, req.Name, req.Email, req.Password, req.PINCode, req.Role, req.IsActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create employee", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusCreated, utils.SuccessResponse("Employee created successfully", employee, nil))
}

func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

	// We need to pass the existing values if not provided, or handle it in service.
	// For simplicity, we just pass the request values to service.
	employee, err := h.employeeService.GetEmployeeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Employee not found", "NOT_FOUND", err.Error()))
		return
	}

	name := req.Name
	if name == "" {
		name = employee.Name
	}
	email := req.Email
	if email == "" {
		email = employee.Email
	}
	role := req.Role
	if role == "" {
		role = employee.Role
	}
	isActive := employee.IsActive
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	updatedEmployee, err := h.employeeService.UpdateEmployee(id, name, email, req.Password, req.PINCode, role, isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update employee", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Employee updated successfully", updatedEmployee, nil))
}

func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	err := h.employeeService.DeleteEmployee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete employee", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Employee deleted successfully", nil, nil))
}
