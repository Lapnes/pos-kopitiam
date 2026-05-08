package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/service"
	"github.com/Lapnes/pos-kopitiam/internal/utils"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	analyticsService service.AnalyticsService
}

func NewAnalyticsHandler(analyticsService service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{analyticsService: analyticsService}
}

func (h *AnalyticsHandler) GetSalesSummary(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()) // Default to start of month
	endDate := now

	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = t
		}
	}

	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Include the whole day
			endDate = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, t.Location())
		}
	}

	summary, err := h.analyticsService.GetSalesSummary(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate sales summary", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Sales summary retrieved", summary, nil))
}

func (h *AnalyticsHandler) GetBestSellers(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	bestSellers, err := h.analyticsService.GetBestSellers(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate best sellers", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Best sellers retrieved", bestSellers, nil))
}

func (h *AnalyticsHandler) GetReturnImpact(c *gin.Context) {
	impact, err := h.analyticsService.GetReturnImpact()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate return impact", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Return impact retrieved", impact, nil))
}
