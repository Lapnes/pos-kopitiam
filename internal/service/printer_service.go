package service

import (
	"fmt"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
)

type PrinterService interface {
	GenerateReceipt(order *models.Order) (string, error)
	GenerateKitchenTicket(order *models.Order) (string, error)
}

type printerService struct{}

func NewPrinterService() PrinterService {
	return &printerService{}
}

func (s *printerService) GenerateReceipt(order *models.Order) (string, error) {
	var details string
	for _, item := range order.OrderDetails {
		if !item.IsVoided {
			details += fmt.Sprintf("  %s x%d @ %.2f = %.2f\n", item.MenuName, item.Quantity, item.Price, item.Subtotal)
		}
	}
	return fmt.Sprintf(
		"=== KOPITIAM RECEIPT ===\nOrder: %s\nDate: %s\n\n%s\nDiscount: %.2f\n========================\n",
		order.OrderNumber, order.CreatedAt.Format("2006-01-02 15:04"), details, order.DiscountAmount,
	), nil
}

func (s *printerService) GenerateKitchenTicket(order *models.Order) (string, error) {
	var items string
	for _, item := range order.OrderDetails {
		if !item.IsVoided {
			items += fmt.Sprintf("  [%s] %s x%d%s\n", item.Station, item.MenuName, item.Quantity, func() string {
				if item.Notes != "" {
					return " (" + item.Notes + ")"
				}
				return ""
			}())
		}
	}
	return fmt.Sprintf(
		"=== KITCHEN TICKET ===\nOrder: %s\nType: %s\n\n%s\n======================\n",
		order.OrderNumber, order.OrderType, items,
	), nil
}
