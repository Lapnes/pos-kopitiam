package repository

import (
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetSalesSummary(startDate, endDate time.Time) (grossRevenue float64, netRevenue float64, totalRefunded float64, totalTransactions int, err error)
	GetBestSellers(limit int) ([]BestSellerResult, error)
	GetReturnImpact() ([]ReturnImpactResult, error)
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
	TotalQty   int     `json:"total_quantity"`
	TotalSales float64 `json:"total_sales"`
}

type ReturnImpactResult struct {
	Reason        string  `json:"reason"`
	ReturnCount   int     `json:"return_count"`
	TotalRefunded float64 `json:"total_refunded"`
}

func (r *analyticsRepository) GetSalesSummary(startDate, endDate time.Time) (float64, float64, float64, int, error) {
	var grossRevenue, totalRefunded float64
	var totalOrders int64

	// Sum Gross Revenue from confirmed/paid orders only (not pending)
	err := r.db.Model(&models.Order{}).
		Where("status IN ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL",
			[]models.OrderStatus{models.OrderConfirmed, models.OrderPaid, models.OrderServed},
			startDate, endDate).
		Select("COALESCE(SUM(total), 0)").
		Row().Scan(&grossRevenue)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Count total transactions (confirmed/paid orders only)
	err = r.db.Model(&models.Order{}).
		Where("status IN ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL",
			[]models.OrderStatus{models.OrderConfirmed, models.OrderPaid, models.OrderServed},
			startDate, endDate).
		Count(&totalOrders).Error
	if err != nil {
		return 0, 0, 0, 0, err
	}

	// Sum Total Refunded — gracefully skip if order_returns table doesn't exist
	var tableExists int64
	r.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'order_returns'").
		Scan(&tableExists)

	if tableExists > 0 {
		err = r.db.Model(&models.OrderReturn{}).
			Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).
			Select("COALESCE(SUM(return_amount), 0)").
			Row().Scan(&totalRefunded)
		if err != nil {
			totalRefunded = 0 // non-fatal, continue with 0
		}
	}

	// Net revenue should never be negative
	netRevenue := grossRevenue - totalRefunded
	if netRevenue < 0 {
		netRevenue = 0
	}

	return grossRevenue, netRevenue, totalRefunded, int(totalOrders), nil
}

func (r *analyticsRepository) GetBestSellers(limit int) ([]BestSellerResult, error) {
	var results []BestSellerResult

	err := r.db.Table("order_details").
		Select("menu_id, menu_name, SUM(quantity) as total_qty, SUM(subtotal) as total_sales").
		Where("deleted_at IS NULL AND is_voided = ?", false).
		Group("menu_id, menu_name").
		Order("total_qty DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

func (r *analyticsRepository) GetReturnImpact() ([]ReturnImpactResult, error) {
	var results []ReturnImpactResult

	// Gracefully handle missing table
	var tableExists int64
	r.db.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = 'order_returns'").
		Scan(&tableExists)

	if tableExists == 0 {
		return results, nil // return empty slice, not error
	}

	err := r.db.Table("order_returns").
		Select("reason, COUNT(id) as return_count, SUM(return_amount) as total_refunded").
		Where("deleted_at IS NULL").
		Group("reason").
		Order("total_refunded DESC").
		Scan(&results).Error

	return results, err
}
