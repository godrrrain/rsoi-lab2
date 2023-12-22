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

	// p, err := psqlDB.GetBooksByLibraryUid(context.Background(), "83556e12-7ce0-48ee-9931-51919ff3c9ee")
	// if err != nil {
	// 	fmt.Printf("Error while getting person: %s", err)
	// } else {
	// 	fmt.Println("Successfully got a person")
	// 	// fmt.Printf("%d: %s, %s, %s, %s \n", p[0].ID, p[0].Library_uid, p[0].Name, p[0].Address, p[0].City)
	// 	fmt.Println(p)
	// 	fmt.Println(len(p))
	// }

	// p, err := psqlDB.GetLibrariesByCity(context.Background(), "Москва")
	// if err != nil {
	// 	fmt.Printf("Error while getting person: %s", err)
	// } else {
	// 	fmt.Println("Successfully got a person")
	// 	// fmt.Printf("%d: %s, %s, %s, %s \n", p[0].ID, p[0].Library_uid, p[0].Name, p[0].Address, p[0].City)
	// 	fmt.Println(p)
	// 	fmt.Println(len(p))
	// }

	handler := handler.NewHandler(psqlDB)

	router := gin.Default()

	router.Use(cors.Default())

	router.GET("/api/v1/libraries/:city", handler.GetLibrariesByCity)
	router.GET("/api/v1/libraries/books/:uid", handler.GetBooksByLibraryUid)

	router.Run()
}
