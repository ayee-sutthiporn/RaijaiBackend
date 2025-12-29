package routes

import (
	"raijai-backend/internal/handlers"
	"raijai-backend/internal/middleware"

	_ "raijai-backend/docs" // Uncommented for Swagger

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB) {
	// Initialize handlers
	userHandler := handlers.NewUserHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	walletHandler := handlers.NewWalletHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)
	debtHandler := handlers.NewDebtHandler(db)
	historyHandler := handlers.NewHistoryLogHandler(db)

	// Swagger Route
	// IMPORTANT: Run 'swag init -g cmd/api/main.go' to generate docs, then uncomment the import above
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// Auth Middleware Group
		// Note: Replace with config values in production
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware("http://localhost:8080/realms/raijai-realm", "raijai-backend"))
		{
			// Users
			protected.POST("/users", userHandler.CreateUser)
			protected.GET("/users/:id", userHandler.GetUser)

			// Categories
			protected.POST("/categories", categoryHandler.CreateCategory)
			protected.GET("/categories", categoryHandler.GetCategories)

			// Wallets
			protected.POST("/wallets", walletHandler.CreateWallet)
			protected.GET("/wallets", walletHandler.GetWallets)

			// Transactions
			protected.POST("/transactions", transactionHandler.CreateTransaction)
			protected.GET("/transactions", transactionHandler.GetTransactions)

			// Debts
			protected.POST("/debts", debtHandler.CreateDebt)
			protected.GET("/debts", debtHandler.GetDebts)

			// History Logs
			protected.GET("/history", historyHandler.GetHistoryLogs)
		}
	}
}
