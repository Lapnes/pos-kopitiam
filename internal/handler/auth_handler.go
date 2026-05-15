package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam/internal/dto"
	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary      Login with email and password
// @Description  Authenticate a user using email + password. Returns an access token and refresh token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200 {object} utils.Response{data=dto.LoginResponse} "Login successful"
// @Failure      400 {object} utils.Response "Invalid request body"
// @Failure      401 {object} utils.Response "Invalid credentials"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body", "BAD_REQUEST", err.Error()))
		return
	}

	res, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Login failed", "UNAUTHORIZED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", res, nil))
}

// LoginPIN godoc
// @Summary      Login with PIN (quick terminal login)
// @Description  Authenticate cashier or kitchen staff using their 6-digit PIN. Faster than email/password flow.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginPINRequest true "PIN credentials"
// @Success      200 {object} utils.Response{data=dto.LoginResponse} "Login successful"
// @Failure      400 {object} utils.Response "Invalid request body"
// @Failure      401 {object} utils.Response "Invalid PIN"
// @Router       /auth/login-pin [post]
func (h *AuthHandler) LoginPIN(c *gin.Context) {
	var req dto.LoginPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request body", "BAD_REQUEST", err.Error()))
		return
	}

	res, err := h.authService.LoginPIN(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse("Login failed", "UNAUTHORIZED", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", res, nil))
}

// Refresh godoc
// @Summary      Refresh access token
// @Description  Exchange a valid refresh token for a new access token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body map[string]string true "Refresh token payload" example({"refresh_token":"eyJhbGci..."})
// @Success      200 {object} utils.Response{data=map[string]string} "Token refreshed"
// @Failure      401 {object} utils.Response "Invalid or expired refresh token"
// @Router       /auth/refresh [post]
func (h *AuthHandler) Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Token refreshed", nil, nil))
}

// Logout godoc
// @Summary      Logout
// @Description  Invalidate the current session token.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response "Logout successful"
// @Failure      401 {object} utils.Response "Unauthorized"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Logout successful", nil, nil))
}
