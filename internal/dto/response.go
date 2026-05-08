package dto

type LoginResponse struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Role     string `json:"role"`
	BranchID string `json:"branch_id"`
}
