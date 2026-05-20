package handler

import (
	"net/http"
	"time"

	"github.com/Lapnes/pos-kopitiam-v3nf/internal/dto"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/models"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/service"
	"github.com/Lapnes/pos-kopitiam-v3nf/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationHandler struct {
	service service.ReservationService
}

func NewReservationHandler(service service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) GetByBranch(c *gin.Context) {
	branchID, err := uuid.Parse(c.Query("branch_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid branch_id", "VALIDATION_ERROR", nil))
		return
	}
	reservations, err := h.service.GetReservationsByBranch(branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to fetch reservations", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Reservations retrieved", reservations, nil))
}

func (h *ReservationHandler) GetByID(c *gin.Context) {
	reservation, err := h.service.GetReservationByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, utils.ErrorResponse("Reservation not found", "NOT_FOUND", nil))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Reservation retrieved", reservation, nil))
}

func (h *ReservationHandler) Create(c *gin.Context) {
	var req dto.ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	branchID, _ := c.Get("branchID")
	branchUUID, _ := uuid.Parse(branchID.(string))
	tableID, _ := uuid.Parse(req.TableID)
	reservationDate, _ := time.Parse("2006-01-02 15:04", req.ReservationDate)
	var customerID *uuid.UUID
	if req.CustomerID != "" {
		id, _ := uuid.Parse(req.CustomerID)
		customerID = &id
	}
	reservation, err := h.service.CreateReservation(branchUUID, tableID, customerID, req.CustomerName, req.CustomerPhone, reservationDate, req.GuestCount, req.DepositAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to create reservation", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusCreated, utils.SuccessResponse("Reservation created", reservation, nil))
}

func (h *ReservationHandler) Update(c *gin.Context) {
	var req dto.ReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, utils.ErrorResponse("Invalid request", "VALIDATION_ERROR", err.Error()))
		return
	}
	reservationDate, _ := time.Parse("2006-01-02 15:04", req.ReservationDate)
	reservation, err := h.service.UpdateReservation(c.Param("id"), req.CustomerName, req.CustomerPhone, reservationDate, req.GuestCount, req.DepositAmount, models.ReservationConfirmed)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to update reservation", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Reservation updated", reservation, nil))
}

func (h *ReservationHandler) Delete(c *gin.Context) {
	if err := h.service.DeleteReservation(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse("Failed to delete reservation", "SERVER_ERROR", err.Error()))
		return
	}
	c.JSON(http.StatusOK, utils.SuccessResponse("Reservation deleted", nil, nil))
}
