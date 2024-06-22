package main

import (
	"dataflow/handlers"
	"dataflow/repo"
	"dataflow/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	repository := repo.NewInMemoryRepository()
	service := services.NewDataService(repository)
	handler := handlers.NewDataHandler(service)

	router := gin.Default()
	router.GET("/data", handler.GetData)
	router.POST("/data", handler.AddData)
	router.POST("/calculate", handler.Calculate)

	log.Println("Server starting on port 8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Could not listen on port 8080: %v\n", err)
	}
}
