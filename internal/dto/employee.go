package dto

type CreateEmployeeRequest struct {
	BranchID string `json:"branch_id" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PINCode  string `json:"pin_code" binding:"required"`
	Role     string `json:"role" binding:"required"`
	IsActive bool   `json:"is_active"`
}

type UpdateEmployeeRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	PINCode  string `json:"pin_code"`
	Role     string `json:"role"`
	IsActive *bool  `json:"is_active"`
}
