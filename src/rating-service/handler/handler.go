package handler

import (
	"context"
	"fmt"
	"net/http"

	"lab2/src/rating-service/storage"

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

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func (h *Handler) GetRating(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	rating, err := h.storage.GetRating(context.Background(), username)

	if err != nil {
		fmt.Printf("failed to get rating %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rating)
}
