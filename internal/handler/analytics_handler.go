package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	service service.AnalyticsService
}

func NewAnalyticsHandler(service service.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func parseDateRange(c *gin.Context) (time.Time, time.Time) {
	startDate, _ := time.Parse("2006-01-02", c.Query("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.Query("end_date"))

	if startDate.IsZero() && endDate.IsZero() {
		dateStr := c.Query("date")
		if dateStr != "" {
			if d, err := time.Parse("2006-01-02", dateStr); err == nil {
				return d, time.Date(d.Year(), d.Month(), d.Day(), 23, 59, 59, 999999999, d.Location())
			}
		}
	}

	if endDate.IsZero() {
		endDate = time.Now()
	} else {
		endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())
	}
	if startDate.IsZero() {
		startDate = endDate.AddDate(0, 0, -30)
		startDate = time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	}
	return startDate, endDate
}

// Legacy endpoints (still work, compute on-demand)
func (h *AnalyticsHandler) SalesSummary(c *gin.Context) {
	startDate, endDate := parseDateRange(c)
	result, err := h.service.GetSalesSummary(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get sales summary", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Sales summary retrieved", result, nil))
}

func (h *AnalyticsHandler) BestSellers(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	startDate, endDate := parseDateRange(c)
	results, err := h.service.GetBestSellers(limit, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get best sellers", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Best sellers retrieved", results, nil))
}

func (h *AnalyticsHandler) ReturnImpact(c *gin.Context) {
	results, err := h.service.GetReturnImpact()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get return impact", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Return impact retrieved", results, nil))
}

func (h *AnalyticsHandler) COGS(c *gin.Context) {
	startDate, _ := time.Parse("2006-01-02", c.Query("start_date"))
	endDate, _ := time.Parse("2006-01-02", c.Query("end_date"))
	if endDate.IsZero() {
		endDate = time.Now()
	}
	if startDate.IsZero() {
		startDate = endDate.AddDate(0, 0, -30)
	}
	cogs, err := h.service.GetCOGS(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get COGS", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("COGS retrieved", gin.H{"cogs": cogs}, nil))
}

// ============================================================
// NEW: Materialized View Endpoints — Instant Dashboard Data
// ============================================================

// DailySales — GET /api/v1/analytics/daily-sales?start_date=2026-05-16&end_date=2026-05-20&branch_id=xxx
func (h *AnalyticsHandler) DailySales(c *gin.Context) {
	startDate, endDate := parseDateRange(c)
	branchID := c.Query("branch_id")
	results, err := h.service.GetDailySales(startDate, endDate, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get daily sales", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Daily sales retrieved", results, nil))
}

// WeeklySales — GET /api/v1/analytics/weekly-sales?week_start=2026-05-12&branch_id=xxx
func (h *AnalyticsHandler) WeeklySales(c *gin.Context) {
	weekStart, _ := time.Parse("2006-01-02", c.Query("week_start"))
	if weekStart.IsZero() {
		// Default to current week start (Monday)
		now := time.Now()
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		weekStart = now.AddDate(0, 0, -(weekday - 1))
	}
	branchID := c.Query("branch_id")
	results, err := h.service.GetWeeklySales(weekStart, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get weekly sales", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Weekly sales retrieved", results, nil))
}

// PaymentSummary — GET /api/v1/analytics/payment-summary?start_date=2026-05-16&end_date=2026-05-20&branch_id=xxx
func (h *AnalyticsHandler) PaymentSummary(c *gin.Context) {
	startDate, endDate := parseDateRange(c)
	branchID := c.Query("branch_id")
	results, err := h.service.GetPaymentSummary(startDate, endDate, branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get payment summary", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Payment summary retrieved", results, nil))
}

// MenuProfitability — GET /api/v1/analytics/menu-profitability
func (h *AnalyticsHandler) MenuProfitability(c *gin.Context) {
	results, err := h.service.GetMenuProfitability()
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get menu profitability", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Menu profitability retrieved", results, nil))
}

// ShiftReconciliation — GET /api/v1/analytics/shift-reconciliation?shift_id=xxx
func (h *AnalyticsHandler) ShiftReconciliation(c *gin.Context) {
	shiftID := c.Query("shift_id")
	results, err := h.service.GetShiftReconciliation(shiftID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get shift reconciliation", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Shift reconciliation retrieved", results, nil))
}

// VoidAndReturnLogs — GET /api/v1/analytics/void-return-logs?days=30
func (h *AnalyticsHandler) VoidAndReturnLogs(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	results, err := h.service.GetVoidAndReturnLogs(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get void/return logs", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Void/return logs retrieved", results, nil))
}

// CashierPerformance — GET /api/v1/analytics/cashier-performance?days=30
func (h *AnalyticsHandler) CashierPerformance(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	results, err := h.service.GetCashierPerformance(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get cashier performance", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Cashier performance retrieved", results, nil))
}

// AuditLogs — GET /api/v1/analytics/audit-logs?limit=100
func (h *AnalyticsHandler) AuditLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	results, err := h.service.GetAuditLogs(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to get audit logs", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Audit logs retrieved", results, nil))
}
