package main

import (
	"context"
	"fmt"

	"lab2/src/library-service/handler"
	"lab2/src/library-service/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	postgresURL := "postgres://program:test@localhost:5432/libraries"
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

	router.GET("/api/v1/libraries", handler.GetLibrariesByCity)
	router.GET("/api/v1/libraries/:uid/books/", handler.GetBooksByLibraryUid)
	router.GET("/api/v1/libraries/:uid/", handler.GetLibraryByUid)
	router.GET("/api/v1/books/:uid/", handler.GetBookInfoByUid)
	router.PUT("/api/v1/books/:uid/count", handler.UpdateBookCount)

	router.Run(":8060")
}
