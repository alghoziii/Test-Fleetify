package main

import (
	"Test_Fleetify/config"
	"Test_Fleetify/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	// Initialize database
	config.ConnectDB()

	// Initialize Gin router
	router := gin.Default()

	// Setup routes
	routes.SetupRoutes(router, config.DB)

	// Start server
	log.Println("Server started on :8080")
	router.Run(":8080")

}
