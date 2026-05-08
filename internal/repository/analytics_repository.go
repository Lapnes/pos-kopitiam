package repository

import (
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetSalesSummary(startDate, endDate time.Time) (grossRevenue float64, netRevenue float64, totalRefunded float64, err error)
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
	MenuID    string  `json:"menu_id"`
	MenuName  string  `json:"menu_name"`
	TotalQty  int     `json:"total_quantity"`
	TotalSales float64 `json:"total_sales"`
}

type ReturnImpactResult struct {
	Reason       string  `json:"reason"`
	ReturnCount  int     `json:"return_count"`
	TotalRefunded float64 `json:"total_refunded"`
}

func (r *analyticsRepository) GetSalesSummary(startDate, endDate time.Time) (float64, float64, float64, error) {
	var grossRevenue, totalRefunded float64

	// Sum Gross Revenue from completed/paid orders
	err := r.db.Model(&models.Order{}).
		Where("status = ? AND created_at BETWEEN ? AND ? AND deleted_at IS NULL", models.OrderPaid, startDate, endDate).
		Select("COALESCE(SUM(total), 0)").
		Row().Scan(&grossRevenue)
	if err != nil {
		return 0, 0, 0, err
	}

	// Sum Total Refunded from order_returns within the timeframe
	err = r.db.Model(&models.OrderReturn{}).
		Where("created_at BETWEEN ? AND ? AND deleted_at IS NULL", startDate, endDate).
		Select("COALESCE(SUM(return_amount), 0)").
		Row().Scan(&totalRefunded)
	if err != nil {
		return 0, 0, 0, err
	}

	netRevenue := grossRevenue - totalRefunded

	return grossRevenue, netRevenue, totalRefunded, nil
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

	err := r.db.Table("order_returns").
		Select("reason, COUNT(id) as return_count, SUM(return_amount) as total_refunded").
		Where("deleted_at IS NULL").
		Group("reason").
		Order("total_refunded DESC").
		Scan(&results).Error

	return results, err
}
