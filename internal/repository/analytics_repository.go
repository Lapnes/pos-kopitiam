package repository

import (
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetSalesSummary(startDate, endDate time.Time) (grossRevenue float64, netRevenue float64, totalRefunded float64, totalTransactions int, err error)
	GetBestSellers(limit int, startDate, endDate time.Time) ([]BestSellerResult, error)
	GetReturnImpact() ([]ReturnImpactResult, error)
	GetCOGS(startDate, endDate time.Time) (float64, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) AnalyticsRepository {
	return &analyticsRepository{db: db}
}

type BestSellerResult struct {
	MenuID     string  `json:"menu_id"`
	MenuName   string  `json:"menu_name"`
	TotalQty   int     `gorm:"column:total_quantity" json:"total_quantity"`
	TotalSales float64 `json:"total_sales"`
}

type ReturnImpactResult struct {
	Reason        string  `json:"reason"`
	ReturnCount   int     `json:"return_count"`
	TotalRefunded float64 `json:"total_refunded"`
}

// GetSalesSummary — FIXED: Replace correlated subquery with JOIN to eliminate N+1
// Also handles case where tables don't exist yet (return zeros instead of error)
func (r *analyticsRepository) GetSalesSummary(startDate, endDate time.Time) (float64, float64, float64, int, error) {
	var grossRevenue, totalRefunded float64
	var totalOrders int64

	// Check if orders table exists
	if !r.db.Migrator().HasTable("orders") {
		return 0, 0, 0, 0, nil
	}

	// FIXED: Use JOIN instead of correlated subquery (1 query, not N+1)
	var grossResult struct{ Gross float64 }
	err := r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(od.subtotal), 0) as gross").
		Joins("LEFT JOIN order_details od ON od.order_id = orders.id AND od.is_voided = 0 AND od.deleted_at IS NULL").
		Where("orders.status IN ? AND orders.created_at BETWEEN ? AND ? AND orders.deleted_at IS NULL",
			[]models.OrderStatus{models.OrderConfirmed, models.OrderPaid, models.OrderServed},
			startDate, endDate).
		Scan(&grossResult).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}
	grossRevenue = grossResult.Gross

	// Count transactions
	err = r.db.Model(&models.Order{}).
		Where("status IN ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL",
			[]models.OrderStatus{models.OrderConfirmed, models.OrderPaid, models.OrderServed},
			startDate, endDate).
		Count(&totalOrders).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Sum refunds - check if table exists
	if r.db.Migrator().HasTable("order_returns") {
		var refundResult struct{ Total float64 }
		err = r.db.Model(&models.OrderReturn{}).
			Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).
			Select("COALESCE(SUM(return_amount), 0) as total").
			Scan(&refundResult).Error
		if err != nil {
			totalRefunded = 0
		} else {
			totalRefunded = refundResult.Total
		}
	}

	netRevenue := grossRevenue - totalRefunded
	if netRevenue < 0 {
		netRevenue = 0
	}

	return grossRevenue, netRevenue, totalRefunded, int(totalOrders), nil
}

// GetBestSellers — FIXED: Use JOIN + GROUP BY instead of per-item lookups
func (r *analyticsRepository) GetBestSellers(limit int, startDate, endDate time.Time) ([]BestSellerResult, error) {
	var results []BestSellerResult

	if !r.db.Migrator().HasTable("order_details") || !r.db.Migrator().HasTable("orders") {
		return []BestSellerResult{}, nil
	}

	err := r.db.Model(&models.OrderDetail{}).
		Select(`
			order_details.menu_id as menu_id,
			order_details.menu_name as menu_name,
			SUM(order_details.quantity) as total_quantity,
			SUM(order_details.subtotal) as total_sales
		`).
		Joins("JOIN orders ON orders.id = order_details.order_id AND orders.deleted_at IS NULL").
		Where("order_details.deleted_at IS NULL AND order_details.is_voided = 0 AND orders.status IN ? AND orders.created_at BETWEEN ? AND ?",
			[]models.OrderStatus{models.OrderConfirmed, models.OrderPaid, models.OrderServed}, startDate, endDate).
		Group("order_details.menu_id, order_details.menu_name").
		Order("total_quantity DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// GetReturnImpact — FIXED: Simple aggregation, no N+1 risk
func (r *analyticsRepository) GetReturnImpact() ([]ReturnImpactResult, error) {
	var results []ReturnImpactResult

	if !r.db.Migrator().HasTable("order_returns") {
		return []ReturnImpactResult{}, nil
	}

	err := r.db.Model(&models.OrderReturn{}).
		Select(`
			reason,
			COUNT(*) as return_count,
			COALESCE(SUM(return_amount), 0) as total_refunded
		`).
		Where("deleted_at IS NULL").
		Group("reason").
		Order("return_count DESC").
		Scan(&results).Error

	return results, err
}

// GetCOGS — FIXED: Use JOIN instead of correlated subquery
func (r *analyticsRepository) GetCOGS(startDate, endDate time.Time) (float64, error) {
	if !r.db.Migrator().HasTable("order_details") || !r.db.Migrator().HasTable("orders") {
		return 0, nil
	}

	var result struct{ Total float64 }

	err := r.db.Model(&models.OrderDetail{}).
		Select("COALESCE(SUM(order_details.cost_price * order_details.quantity), 0) as total").
		Joins("JOIN orders ON orders.id = order_details.order_id AND orders.deleted_at IS NULL").
		Where("order_details.deleted_at IS NULL AND order_details.is_voided = 0 AND orders.status IN ? AND orders.created_at BETWEEN ? AND ?",
			[]models.OrderStatus{models.OrderConfirmed, models.OrderPaid, models.OrderServed},
			startDate, endDate).
		Scan(&result).Error

	return result.Total, err
}
