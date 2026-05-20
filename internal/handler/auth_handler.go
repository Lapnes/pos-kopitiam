package handler

import (
	"net/http"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Login godoc
// @Summary Login with email/password
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginRequest true "Login credentials"
// @Success 200 {object} utils.Response
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	resp, err := h.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error(), "UNAUTHORIZED", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", resp, nil))
}

// LoginPIN godoc
// @Summary Login with PIN
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.LoginPINRequest true "PIN"
// @Success 200 {object} utils.Response
// @Router /auth/login-pin [post]
func (h *AuthHandler) LoginPIN(c *gin.Context) {
	var req dto.LoginPINRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	resp, err := h.authService.LoginPIN(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, utils.ErrorResponse(err.Error(), "UNAUTHORIZED", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Login successful", resp, nil))
}
