package main

import (
	"lab2/src/gateway-service/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	handler := handler.NewHandler()

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/v1/libraries", handler.GetLibrariesByCity)
	router.GET("/api/v1/libraries/:uid/books/", handler.GetBooksByLibraryUid)
	router.GET("/api/v1/rating/", handler.GetRating)
	router.GET("/api/v1/reservations", handler.GetReservations)
	router.POST("/api/v1/reservations", handler.CreateReservation)
	router.POST("/api/v1/reservations/:uid/return", handler.ReturnBook)

	router.GET("/manage/health", handler.GetHealth)

	router.Run()
}
