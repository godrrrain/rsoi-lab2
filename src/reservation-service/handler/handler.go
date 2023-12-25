package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"lab2/src/reservation-service/storage"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type Handler struct {
	storage storage.Storage
}

type RequestCreateReservation struct {
	BookUid    string `json:"bookUid"`
	LibraryUid string `json:"libraryUid"`
	TillDate   string `json:"tillDate"`
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetReservations(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	reservations, err := h.storage.GetReservations(context.Background(), username)

	if err != nil {
		fmt.Printf("failed to get reservations %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reservations)
}

func (h *Handler) GetReservationByUid(c *gin.Context) {

	reservation, err := h.storage.GetReservationByUid(context.Background(), c.Param("uid"))

	if err != nil {
		fmt.Printf("failed to get reservation %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func (h *Handler) GetRentedReservationAmount(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	reservationAmount, err := h.storage.GetRentedReservationAmount(context.Background(), username)

	if err != nil {
		fmt.Printf("failed to get reservation amount %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reservationAmount)
}

func (h *Handler) CreateReservation(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	var reqCrRes RequestCreateReservation

	err := json.NewDecoder(c.Request.Body).Decode(&reqCrRes)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	reservation, err := h.storage.CreateReservation(context.Background(), username, reqCrRes.BookUid, reqCrRes.LibraryUid, reqCrRes.TillDate)

	if err != nil {
		fmt.Printf("failed to create reservations %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, reservation)
}

func (h *Handler) UpdateReservationStatus(c *gin.Context) {

	err := h.storage.UpdateReservationStatus(context.Background(), c.Param("uid"), c.Query("condition"))

	if err != nil {
		fmt.Printf("failed to update reservation %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "status updated",
	})
}
