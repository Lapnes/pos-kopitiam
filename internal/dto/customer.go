package dto

type CreateCustomerRequest struct {
	Name  string `json:"name" binding:"required"`
	Phone string `json:"phone" binding:"required"`
	Email string `json:"email" binding:"omitempty,email"`
}

type UpdateCustomerRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email" binding:"omitempty,email"`
}

type CustomerResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Phone      string  `json:"phone"`
	Email      string  `json:"email"`
	VisitCount int     `json:"visit_count"`
	TotalSpent float64 `json:"total_spent"`
}