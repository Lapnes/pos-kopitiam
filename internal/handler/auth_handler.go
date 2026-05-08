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

func (h *AuthHandler) Refresh(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Token refreshed", nil, nil))
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.JSON(http.StatusOK, utils.SuccessResponse("Logout successful", nil, nil))
}
