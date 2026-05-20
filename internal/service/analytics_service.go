package service

import (
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"gorm.io/gorm"
)

type AnalyticsService interface {
	GetSalesSummary(startDate, endDate time.Time) (map[string]interface{}, error)
	GetBestSellers(limit int, startDate, endDate time.Time) ([]map[string]interface{}, error)
	GetReturnImpact() ([]map[string]interface{}, error)
	GetCOGS(startDate, endDate time.Time) (float64, error)

	// New: Materialized view queries for instant dashboard data
	GetDailySales(startDate, endDate time.Time, branchID string) ([]map[string]interface{}, error)
	GetWeeklySales(weekStart time.Time, branchID string) ([]map[string]interface{}, error)
	GetPaymentSummary(startDate, endDate time.Time, branchID string) ([]map[string]interface{}, error)
	GetMenuProfitability() ([]map[string]interface{}, error)
	GetShiftReconciliation(shiftID string) ([]map[string]interface{}, error)
	GetVoidAndReturnLogs(days int) ([]map[string]interface{}, error)
	GetCashierPerformance(days int) ([]map[string]interface{}, error)
	GetAuditLogs(limit int) ([]models.AuditLog, error)
}

type analyticsService struct {
	db            *gorm.DB
	analyticsRepo repository.AnalyticsRepository
}

func NewAnalyticsService(db *gorm.DB, analyticsRepo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsService{db: db, analyticsRepo: analyticsRepo}
}

// Legacy methods (still work, but slower)
func (s *analyticsService) GetSalesSummary(startDate, endDate time.Time) (map[string]interface{}, error) {
	gross, net, refunded, transactions, err := s.analyticsRepo.GetSalesSummary(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate previous period of identical duration
	duration := endDate.Sub(startDate)
	prevStartDate := startDate.Add(-duration)
	prevEndDate := startDate.Add(-time.Second)

	prevGross, prevNet, _, prevTransactions, err := s.analyticsRepo.GetSalesSummary(prevStartDate, prevEndDate)
	if err != nil {
		prevGross, prevNet, prevTransactions = 0.0, 0.0, 0
	}

	// Gross Change
	var grossChange float64
	if prevGross > 0 {
		grossChange = ((gross - prevGross) / prevGross) * 100
	} else if gross > 0 {
		grossChange = 100
	}

	// Net Change
	var netChange float64
	if prevNet > 0 {
		netChange = ((net - prevNet) / prevNet) * 100
	} else if net > 0 {
		netChange = 100
	}

	// Transactions Change
	var transactionsChange float64
	if prevTransactions > 0 {
		transactionsChange = float64(transactions-prevTransactions) / float64(prevTransactions) * 100
	} else if transactions > 0 {
		transactionsChange = 100
	}

	// AOV Change
	aov := float64(0)
	if transactions > 0 {
		aov = net / float64(transactions)
	}
	prevAov := float64(0)
	if prevTransactions > 0 {
		prevAov = prevNet / float64(prevTransactions)
	}
	var aovChange float64
	if prevAov > 0 {
		aovChange = ((aov - prevAov) / prevAov) * 100
	} else if aov > 0 {
		aovChange = 100
	}

	return map[string]interface{}{
		"gross_revenue":             gross,
		"net_revenue":               net,
		"total_refunded":            refunded,
		"total_transactions":        transactions,
		"gross_revenue_change":      grossChange,
		"net_revenue_change":        netChange,
		"total_transactions_change": transactionsChange,
		"avg_order_value_change":    aovChange,
	}, nil
}

func (s *analyticsService) GetBestSellers(limit int, startDate, endDate time.Time) ([]map[string]interface{}, error) {
	results, err := s.analyticsRepo.GetBestSellers(limit, startDate, endDate)
	if err != nil {
		return nil, err
	}
	output := make([]map[string]interface{}, len(results))
	for i, r := range results {
		output[i] = map[string]interface{}{
			"menu_id":        r.MenuID,
			"menu_name":      r.MenuName,
			"total_quantity": r.TotalQty,
			"total_sales":    r.TotalSales,
		}
	}
	return output, nil
}

func (s *analyticsService) GetReturnImpact() ([]map[string]interface{}, error) {
	results, err := s.analyticsRepo.GetReturnImpact()
	if err != nil {
		return nil, err
	}
	output := make([]map[string]interface{}, len(results))
	for i, r := range results {
		output[i] = map[string]interface{}{
			"reason":         r.Reason,
			"return_count":   r.ReturnCount,
			"total_refunded": r.TotalRefunded,
		}
	}
	return output, nil
}

func (s *analyticsService) GetCOGS(startDate, endDate time.Time) (float64, error) {
	return s.analyticsRepo.GetCOGS(startDate, endDate)
}

// ============================================================
// NEW: Materialized View Queries — Instant Dashboard Data
// ============================================================

// GetDailySales — reads from v_daily_sales_summary (refreshed every 5 min)
func (s *analyticsService) GetDailySales(startDate, endDate time.Time, branchID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := s.db.Table("v_daily_sales_summary").
		Where("sale_date BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}
	err := query.Order("sale_date ASC").Find(&results).Error
	return results, err
}

// GetWeeklySales — reads from v_weekly_sales_summary
func (s *analyticsService) GetWeeklySales(weekStart time.Time, branchID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := s.db.Table("v_weekly_sales_summary").
		Where("week_start = ?", weekStart.Format("2006-01-02"))
	if branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}
	err := query.Find(&results).Error
	return results, err
}

// GetPaymentSummary — reads from v_payment_method_summary
func (s *analyticsService) GetPaymentSummary(startDate, endDate time.Time, branchID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := s.db.Table("v_payment_method_summary").
		Select("payment_method, SUM(total_amount) as total_amount").
		Where("payment_date BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if branchID != "" {
		query = query.Where("branch_id = ?", branchID)
	}
	err := query.Group("payment_method").Order("total_amount DESC").Find(&results).Error
	return results, err
}

// GetMenuProfitability — reads from v_menu_profitability_live
func (s *analyticsService) GetMenuProfitability() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := s.db.Table("v_menu_profitability_live").
		Order("net_profit DESC").
		Find(&results).Error
	return results, err
}

// GetBestSellersFromView — reads from v_best_selling_items (instant, no GROUP BY)
func (s *analyticsService) GetBestSellersFromView(limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := s.db.Table("v_best_selling_items").
		Order("total_quantity DESC").
		Limit(limit).
		Find(&results).Error
	return results, err
}

// GetShiftReconciliation — reads from v_shift_reconciliation
func (s *analyticsService) GetShiftReconciliation(shiftID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := s.db.Table("v_shift_reconciliation")
	if shiftID != "" {
		query = query.Where("shift_id = ?", shiftID)
	}
	err := query.Order("opening_time DESC").Find(&results).Error
	return results, err
}

// GetVoidAndReturnLogs — reads from v_void_and_return_logs
func (s *analyticsService) GetVoidAndReturnLogs(days int) ([]map[string]interface{}, error) {
	if days <= 0 || days > 90 {
		days = 30 // Default + safety limit
	}
	var results []map[string]interface{}
	err := s.db.Table("v_void_and_return_logs").
		Where("event_date >= DATE_SUB(CURRENT_DATE, INTERVAL ? DAY)", days).
		Order("event_date DESC, amount DESC").
		Find(&results).Error
	return results, err
}

// GetCashierPerformance — reads from v_cashier_performance
func (s *analyticsService) GetCashierPerformance(days int) ([]map[string]interface{}, error) {
	if days <= 0 || days > 90 {
		days = 30
	}
	var results []map[string]interface{}
	err := s.db.Table("v_cashier_performance").
		Order("total_sales DESC").
		Find(&results).Error
	return results, err
}

// GetAuditLogs — reads from audit_logs
func (s *analyticsService) GetAuditLogs(limit int) ([]models.AuditLog, error) {
	if limit <= 0 {
		limit = 100
	}
	var logs []models.AuditLog
	err := s.db.Order("timestamp DESC").Limit(limit).Find(&logs).Error
	return logs, err
}
