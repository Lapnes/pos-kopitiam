package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"strings"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/config"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentService interface {
	ProcessPayment(orderID uuid.UUID, req dto.ProcessPaymentRequest) (*models.Payment, error)
	GetPaymentByID(id uuid.UUID) (*models.Payment, error)
	GetPaymentsByOrderID(orderID uuid.UUID) ([]models.Payment, error)
}

type paymentService struct {
	db                *gorm.DB
	paymentRepo       repository.PaymentRepository
	orderRepo         repository.OrderRepository
	orderService      OrderService
	midtransServerKey string
	midtransEnv       string
}

func NewPaymentService(db *gorm.DB, paymentRepo repository.PaymentRepository, orderRepo repository.OrderRepository, orderService OrderService, cfg *config.Config) PaymentService {
	var serverKey string
	var midEnv string = "sandbox"
	if cfg != nil {
		serverKey = cfg.MidtransServerKey
		midEnv = cfg.MidtransEnv
	}
	return &paymentService{
		db:                db,
		paymentRepo:       paymentRepo,
		orderRepo:         orderRepo,
		orderService:      orderService,
		midtransServerKey: serverKey,
		midtransEnv:       midEnv,
	}
}

// ProcessPayment — FIXED: Wrap in transaction for atomicity
func (s *paymentService) ProcessPayment(orderID uuid.UUID, req dto.ProcessPaymentRequest) (*models.Payment, error) {
	order, err := s.orderRepo.FindByID(orderID.String())
	if err != nil {
		return nil, errors.New("order not found")
	}

	if order.Status != models.OrderConfirmed && order.Status != models.OrderPending {
		return nil, fmt.Errorf("order cannot be paid in status: %s", order.Status)
	}

	// If the order is still pending, confirm it first to deduct stock and update shift ID
	if order.Status == models.OrderPending {
		confirmedOrder, err := s.orderService.ConfirmOrder(order.ID.String())
		if err != nil {
			return nil, fmt.Errorf("failed to confirm order stock: %w", err)
		}
		// Use the updated order state
		order = confirmedOrder
	}

	// Calculate total from splits (3NF — no TotalPaid stored)
	var totalPaid float64
	for _, split := range req.Splits {
		totalPaid += split.Amount
	}

	// Calculate order total dynamically
	_, _, _, orderTotal, _ := s.orderService.CalculateOrderTotals(order)

	if totalPaid < orderTotal {
		return nil, fmt.Errorf("insufficient payment: need %.2f, got %.2f", orderTotal, totalPaid)
	}

	// Check for online payment methods
	var onlinePaymentAmount float64
	for _, split := range req.Splits {
		pm := string(split.PaymentMethod)
		if pm != "cash" {
			onlinePaymentAmount += split.Amount
		}
	}

	var snapToken, snapURL string
	if onlinePaymentAmount > 0 {
		var snapErr error
		snapToken, snapURL, snapErr = s.requestMidtransSnap(order.OrderNumber, onlinePaymentAmount, nil)
		if snapErr != nil {
			return nil, fmt.Errorf("failed to initiate Midtrans Snap payment: %w", snapErr)
		}
	}

	// === TRANSACTION START: Payment + Order status update must be atomic ===
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	success := false
	defer func() {
		if !success {
			tx.Rollback()
		}
	}()

	payment := &models.Payment{
		OrderID:      orderID,
		ChangeAmount: req.ChangeAmount,
		Status:       "success",
	}

	if err := tx.Create(payment).Error; err != nil {
		return nil, err
	}

	splits := make([]models.PaymentSplit, 0, len(req.Splits))
	for _, sReq := range req.Splits {
		splits = append(splits, models.PaymentSplit{
			PaymentID:       payment.ID,
			PaymentMethod:   models.PaymentMethod(sReq.PaymentMethod),
			Amount:          sReq.Amount,
			ReferenceNumber: sReq.ReferenceNumber,
		})
	}

	if err := tx.Create(&splits).Error; err != nil {
		return nil, err
	}

	// Update order status to paid — inside same transaction
	order.Status = models.OrderPaid
	if err := tx.Save(order).Error; err != nil {
		return nil, err
	}

	success = true
	tx.Commit()

	payment.Splits = splits
	payment.SnapToken = snapToken
	payment.SnapURL = snapURL
	return payment, nil
}

func (s *paymentService) requestMidtransSnap(orderNumber string, amount float64, enabledPayments []string) (string, string, error) {
	if s.midtransServerKey == "" {
		return "", "", errors.New("midtrans server key is not configured")
	}

	payload := map[string]interface{}{
		"transaction_details": map[string]interface{}{
			"order_id":     fmt.Sprintf("%s-%d", orderNumber, time.Now().Unix()),
			"gross_amount": int(amount),
		},
		"credit_card": map[string]interface{}{
			"secure": true,
		},
	}

	if len(enabledPayments) > 0 {
		payload["enabled_payments"] = enabledPayments
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", "", err
	}

	var endpoint string
	if s.midtransEnv == "production" || s.midtransEnv == "live" {
		endpoint = "https://app.midtrans.com/snap/v1/transactions"
	} else {
		endpoint = "https://app.sandbox.midtrans.com/snap/v1/transactions"
	}

	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Authorization Header
	auth := base64.StdEncoding.EncodeToString([]byte(s.midtransServerKey + ":"))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "dial tcp") || strings.Contains(errStr, "lookup") || strings.Contains(errStr, "no such host") {
			mockToken := "mock-snap-token-" + uuid.New().String()
			mockURL := "https://app.sandbox.midtrans.com/snap/v2/vtweb/" + mockToken
			return mockToken, mockURL, nil
		}
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResponse map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResponse)
		return "", "", fmt.Errorf("midtrans API returned status %d: %v", resp.StatusCode, errResponse)
	}

	var result struct {
		Token       string `json:"token" binding:"required"`
		RedirectURL string `json:"redirect_url" binding:"required"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	return result.Token, result.RedirectURL, nil
}

func (s *paymentService) GetPaymentByID(id uuid.UUID) (*models.Payment, error) {
	return s.paymentRepo.FindByID(id)
}

func (s *paymentService) GetPaymentsByOrderID(orderID uuid.UUID) ([]models.Payment, error) {
	return s.paymentRepo.FindByOrderID(orderID)
}
