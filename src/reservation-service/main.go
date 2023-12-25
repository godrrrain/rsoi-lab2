package main

import (
	"context"
	"fmt"

	"lab2/src/reservation-service/handler"
	"lab2/src/reservation-service/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	postgresURL := "postgres://program:test@localhost:5432/reservations"
	psqlDB, err := storage.NewPgStorage(context.Background(), postgresURL)
	if err != nil {
		fmt.Printf("Postgresql init: %s", err)
	} else {
		fmt.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	handler := handler.NewHandler(psqlDB)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/v1/reservations", handler.GetReservations)
	router.GET("/api/v1/reservations/info/:uid", handler.GetReservationByUid)
	router.GET("/api/v1/reservations/amount", handler.GetRentedReservationAmount)
	router.POST("/api/v1/reservations", handler.CreateReservation)
	router.PATCH("/api/v1/reservations/:uid", handler.UpdateReservationStatus)

	router.Run()
}
