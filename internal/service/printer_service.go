package service

import (
	"errors"
	"time"

	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/google/uuid"
)

type PrinterService interface {
	GenerateReceipt(orderID uuid.UUID) (map[string]interface{}, error)
	GenerateKitchenTicket(orderID uuid.UUID) (map[string]interface{}, error)
}

type printerService struct {
	orderRepo repository.OrderRepository
	userRepo  repository.UserRepository
}

func NewPrinterService(orderRepo repository.OrderRepository, userRepo repository.UserRepository) PrinterService {
	return &printerService{orderRepo: orderRepo, userRepo: userRepo}
}

func (s *printerService) GenerateReceipt(orderID uuid.UUID) (map[string]interface{}, error) {
	order, err := s.orderRepo.FindByID(orderID.String())
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}

	cashierName := "Unknown Cashier"
	if cashier, err := s.userRepo.FindByID(order.EmployeeID.String()); err == nil {
		cashierName = cashier.Name
	}

	var items []map[string]interface{}
	for _, detail := range order.OrderDetails {
		items = append(items, map[string]interface{}{
			"menu_name":  detail.MenuName,
			"quantity":   detail.Quantity,
			"unit_price": detail.Price,
			"subtotal":   detail.Subtotal,
		})
	}

	var payments []map[string]interface{}
	for _, payment := range order.Payments {
		payments = append(payments, map[string]interface{}{
			"method": string(payment.PaymentMethod),
			"amount": payment.Amount,
		})
	}

	payload := map[string]interface{}{
		"type":         "receipt",
		"store_header": map[string]interface{}{
			"name":    "KopiTiam Main Branch",
			"address": "Jl. Sudirman No. 1, Jakarta",
			"phone":   "+628123456789",
		},
		"date":         time.Now().Format("2006-01-02 15:04:05"),
		"order_number": order.OrderNumber,
		"order_type":   string(order.OrderType),
		"cashier":      cashierName,
		"items":        items,
		"totals": map[string]interface{}{
			"subtotal":       order.Subtotal,
			"tax":            order.TaxAmount,
			"service_charge": order.ServiceCharge,
			"discount":       order.DiscountAmount,
			"grand_total":    order.Total,
		},
		"payments": payments,
		"footer":   "Thank you for visiting KopiTiam! See you again.",
	}

	return payload, nil
}

func (s *printerService) GenerateKitchenTicket(orderID uuid.UUID) (map[string]interface{}, error) {
	order, err := s.orderRepo.FindByID(orderID.String())
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}

	cashierName := "Unknown Cashier"
	if cashier, err := s.userRepo.FindByID(order.EmployeeID.String()); err == nil {
		cashierName = cashier.Name
	}

	var tableNumber string
	if order.TableID != nil {
		tableNumber = order.TableID.String() // Simplified, should query table name
	} else {
		tableNumber = "Takeaway/Delivery"
	}

	var items []map[string]interface{}
	for _, detail := range order.OrderDetails {
		items = append(items, map[string]interface{}{
			"menu_name": detail.MenuName,
			"quantity":  detail.Quantity,
			"station":   string(detail.Station),
			"notes":     detail.Notes,
		})
	}

	payload := map[string]interface{}{
		"type":         "kitchen_ticket",
		"order_type":   string(order.OrderType),
		"table":        tableNumber,
		"order_number": order.OrderNumber,
		"time":         time.Now().Format("15:04:05"),
		"cashier":      cashierName,
		"items":        items,
	}

	return payload, nil
}
