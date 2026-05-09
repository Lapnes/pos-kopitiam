package service

import (
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/repository"
)

type AnalyticsService interface {
	GetSalesSummary(startDate, endDate time.Time) (map[string]interface{}, error)
	GetBestSellers(limit int) ([]repository.BestSellerResult, error)
	GetReturnImpact() ([]repository.ReturnImpactResult, error)
}

type analyticsService struct {
	analyticsRepo repository.AnalyticsRepository
}

func NewAnalyticsService(analyticsRepo repository.AnalyticsRepository) AnalyticsService {
	return &analyticsService{analyticsRepo: analyticsRepo}
}

func (s *analyticsService) GetSalesSummary(startDate, endDate time.Time) (map[string]interface{}, error) {
	gross, net, refunded, transactions, err := s.analyticsRepo.GetSalesSummary(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"start_date":         startDate.Format(time.RFC3339),
		"end_date":           endDate.Format(time.RFC3339),
		"gross_revenue":      gross,
		"net_revenue":        net,
		"total_revenue":      net,
		"total_refunded":     refunded,
		"total_transactions": transactions,
	}, nil
}

func (s *analyticsService) GetBestSellers(limit int) ([]repository.BestSellerResult, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.analyticsRepo.GetBestSellers(limit)
}

func (s *analyticsService) GetReturnImpact() ([]repository.ReturnImpactResult, error) {
	return s.analyticsRepo.GetReturnImpact()
}
