package main

import (
	"Test_Fleetify/config"
	"Test_Fleetify/routes"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"log"
	"time"
)

func main() {
	// Initialize database
	config.ConnectDB()

	// Initialize Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Setup routes
	routes.SetupRoutes(router, config.DB)

	// Start server
	log.Println("Server started on :8080")
	router.Run(":8080")

}
