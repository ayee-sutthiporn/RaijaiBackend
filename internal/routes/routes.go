package routes

import (
	"raijai-backend/internal/config"
	"raijai-backend/internal/handlers"
	"raijai-backend/internal/middleware"

	_ "raijai-backend/docs" // Uncommented for Swagger

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)
 
func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg *config.Config) {
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
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(cfg.KeycloakIssuer, cfg.KeycloakClientID))
		{
			// Auth / Profile
			protected.GET("/auth/me", userHandler.GetMe)

			// Users
			protected.POST("/users", userHandler.CreateUser)
			protected.GET("/users/:id", userHandler.GetUser)

			// Categories
			protected.POST("/categories", categoryHandler.CreateCategory)
			protected.GET("/categories", categoryHandler.GetCategories)
			protected.PUT("/categories/:id", categoryHandler.UpdateCategory)
			protected.DELETE("/categories/:id", categoryHandler.DeleteCategory)

			// Wallets
			protected.POST("/wallets", walletHandler.CreateWallet)
			protected.GET("/wallets", walletHandler.GetWallets)
			protected.PUT("/wallets/:id", walletHandler.UpdateWallet)
			protected.DELETE("/wallets/:id", walletHandler.DeleteWallet)

			// Transactions
			protected.POST("/transactions", transactionHandler.CreateTransaction)
			protected.GET("/transactions", transactionHandler.GetTransactions)
			protected.PUT("/transactions/:id", transactionHandler.UpdateTransaction)
			protected.DELETE("/transactions/:id", transactionHandler.DeleteTransaction)

			// Debts
			protected.POST("/debts", debtHandler.CreateDebt)
			protected.GET("/debts", debtHandler.GetDebts)
			protected.PUT("/debts/:id", debtHandler.UpdateDebt)
			protected.DELETE("/debts/:id", debtHandler.DeleteDebt)
			protected.POST("/debts/:id/payment", debtHandler.MakePayment)

			// History Logs
			protected.GET("/history", historyHandler.GetHistoryLogs)
		}
	}
}
