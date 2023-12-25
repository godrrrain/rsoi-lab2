package handler

import (
	"context"
	"fmt"
	"net/http"

	"lab2/src/library-service/storage"

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

func (h *Handler) GetLibrariesByCity(c *gin.Context) {

	libraries, err := h.storage.GetLibrariesByCity(context.Background(), c.Query("city"))

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, libraries)
}

func (h *Handler) GetBooksByLibraryUid(c *gin.Context) {

	libraries, err := h.storage.GetBooksByLibraryUid(context.Background(), c.Param("uid"))

	if err != nil {
		fmt.Printf("failed to get libraries %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, libraries)
}
