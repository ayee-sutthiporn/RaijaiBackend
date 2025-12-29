package main

import (
	"log"
	"raijai-backend/internal/config"
	"raijai-backend/internal/models"
	"raijai-backend/internal/routes"

	"github.com/gin-gonic/gin"
)

// @title RaiJai API
// @version 1.0
// @description RESTful API for RaiJai Application
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Load configuration and connect to database
	cfg := config.LoadConfig()
	db := config.ConnectDB(cfg)

	// Auto Migrate Database
	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("Failed to duplicate database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup routes
	routes.SetupRoutes(r, db)

	// Start server
	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
