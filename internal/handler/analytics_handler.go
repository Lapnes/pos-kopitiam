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

// GetSalesSummary godoc
// @Summary      Get sales summary for a date range
// @Description  Returns gross revenue, net revenue, total refunded, and total transaction count for confirmed/paid/served orders within the given date range. Defaults to the current calendar month if no dates are provided.
// @Tags         Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        start_date query string false "Start date in YYYY-MM-DD format (default: first day of current month)"
// @Param        end_date query string false "End date in YYYY-MM-DD format (default: today)"
// @Success      200 {object} utils.Response{data=map[string]interface{}} "Sales summary"
// @Failure      500 {object} utils.Response "Failed to generate sales summary"
// @Router       /analytics/sales-summary [get]
func (h *AnalyticsHandler) GetSalesSummary(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endDate := now

	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = t
		}
	}

	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
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

// GetBestSellers godoc
// @Summary      Get top-selling menu items
// @Description  Returns the top N menu items ranked by total quantity sold (non-voided order details only). Defaults to top 10 if limit is not specified.
// @Tags         Analytics
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Number of top items to return (default: 10, min: 1)"
// @Success      200 {object} utils.Response{data=[]repository.BestSellerResult} "Best sellers list"
// @Failure      500 {object} utils.Response "Failed to generate best sellers"
// @Router       /analytics/best-sellers [get]
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

// GetReturnImpact godoc
// @Summary      Get return/refund impact report
// @Description  Returns a breakdown of order returns grouped by reason, showing return_count and total_refunded amount for each reason. Useful for identifying patterns in customer complaints.
// @Tags         Analytics
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} utils.Response{data=[]repository.ReturnImpactResult} "Return impact report"
// @Failure      500 {object} utils.Response "Failed to generate return impact"
// @Router       /analytics/return-impact [get]
func (h *AnalyticsHandler) GetReturnImpact(c *gin.Context) {
	impact, err := h.analyticsService.GetReturnImpact()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to generate return impact", "INTERNAL_ERROR", err.Error()))
		return
	}

	c.JSON(http.StatusOK, utils.SuccessResponse("Return impact retrieved", impact, nil))
}
