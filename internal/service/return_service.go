package service

import (
	"errors"

	"github.com/Lapnes/pos-kopitiam/internal/models"
	"github.com/Lapnes/pos-kopitiam/internal/repository"
	"github.com/google/uuid"
)

type ReturnService interface {
	ProcessReturn(orderID, processedBy uuid.UUID, returnAmount float64, reason string) (*models.OrderReturn, error)
}

type returnService struct {
	returnRepo repository.ReturnRepository
	orderRepo  repository.OrderRepository
}

func NewReturnService(returnRepo repository.ReturnRepository, orderRepo repository.OrderRepository) ReturnService {
	return &returnService{
		returnRepo: returnRepo,
		orderRepo:  orderRepo,
	}
}

func (s *returnService) ProcessReturn(orderID, processedBy uuid.UUID, returnAmount float64, reason string) (*models.OrderReturn, error) {
	order, err := s.orderRepo.FindByID(orderID.String())
	if err != nil {
		return nil, errors.New("order not found: " + err.Error())
	}

	// Rule: Maximum returnable amount is 80% of the total price.
	maxAllowedReturn := order.Total * 0.80

	if returnAmount > maxAllowedReturn {
		return nil, errors.New("return amount exceeds the maximum 80% limit")
	}

	orderReturn := &models.OrderReturn{
		OrderID:        orderID,
		ProcessedBy:    processedBy,
		Reason:         reason,
		OriginalAmount: order.Total,
		ReturnAmount:   returnAmount,
	}

	// Save the return record
	if err := s.returnRepo.Create(orderReturn); err != nil {
		return nil, errors.New("failed to process return: " + err.Error())
	}

	// Flag the order as having returns
	order.HasReturns = true
	if err := s.orderRepo.Update(order); err != nil {
		return nil, errors.New("return processed but failed to update order flag: " + err.Error())
	}

	return orderReturn, nil
}
