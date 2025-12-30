package main

import (
	"log"
	"raijai-backend/internal/config"
	"raijai-backend/internal/models"
	"raijai-backend/internal/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title RaiJai API
// @version 1.0
// @description RESTful API for RaiJai Application
// @host raijai-api.sutthiporn.dev
// @schemes https
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

	// CORS Configuration
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "https://raijai.sutthiporn.dev"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup routes
	routes.SetupRoutes(r, db, cfg)

	// Start server
	log.Println("Server starting on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
