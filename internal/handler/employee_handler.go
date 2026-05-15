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

// GetEmployees godoc
// @Summary      List all employees
// @Description  Returns all employees across all branches. Manager role required.
// @Tags         Employees
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response{data=[]models.Employee} "List of employees"
// @Failure      500 {object} utils.Response "Failed to fetch employees"
// @Router       /employees [get]
func (h *EmployeeHandler) GetEmployees(c *gin.Context) {
	employees, err := h.employeeService.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch employees", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Employees retrieved successfully", employees, nil))
}

// CreateEmployee godoc
// @Summary      Create a new employee
// @Description  Creates a new employee with one of the allowed roles: manager | cashier | kitchen. The email must be unique. PIN must be 6 digits and unique within the branch.
// @Tags         Employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body CreateEmployeeRequest true "Employee creation payload"
// @Success      201 {object} utils.Response{data=models.Employee} "Employee created successfully"
// @Failure      400 {object} utils.Response "Invalid payload or invalid branch UUID"
// @Failure      500 {object} utils.Response "Failed to create employee"
// @Router       /employees [post]
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

// UpdateEmployee godoc
// @Summary      Update an employee
// @Description  Partially updates an employee record. Fields left empty retain their current values. Password and PIN are updated only if provided (non-empty strings).
// @Tags         Employees
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Employee UUID"
// @Param        request body UpdateEmployeeRequest true "Fields to update (all optional)"
// @Success      200 {object} utils.Response{data=models.Employee} "Employee updated successfully"
// @Failure      400 {object} utils.Response "Invalid payload"
// @Failure      404 {object} utils.Response "Employee not found"
// @Failure      500 {object} utils.Response "Failed to update employee"
// @Router       /employees/{id} [put]
func (h *EmployeeHandler) UpdateEmployee(c *gin.Context) {
	id := c.Param("id")
	var req UpdateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request payload", "BAD_REQUEST", err.Error()))
		return
	}

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

// DeleteEmployee godoc
// @Summary      Delete an employee
// @Description  Soft-deletes an employee by UUID. The employee record is retained in the database (deleted_at is set).
// @Tags         Employees
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Employee UUID"
// @Success      200 {object} utils.Response "Employee deleted successfully"
// @Failure      500 {object} utils.Response "Failed to delete employee"
// @Router       /employees/{id} [delete]
func (h *EmployeeHandler) DeleteEmployee(c *gin.Context) {
	id := c.Param("id")

	err := h.employeeService.DeleteEmployee(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete employee", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Employee deleted successfully", nil, nil))
}
