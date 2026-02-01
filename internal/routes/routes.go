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
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)
	categoryHandler := handlers.NewCategoryHandler(db)
	walletHandler := handlers.NewWalletHandler(db)
	transactionHandler := handlers.NewTransactionHandler(db)
	debtHandler := handlers.NewDebtHandler(db)
	historyHandler := handlers.NewHistoryLogHandler(db)
	reportHandler := handlers.NewReportHandler(db)
	systemHandler := handlers.NewSystemHandler(db)
	bookHandler := handlers.NewBookHandler(db)

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		// Public Auth Routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh-token", authHandler.RefreshToken)
		}

		// Protected Middleware Group
		protected := api.Group("/")
		// Note: AuthMiddleware now accepts 'mock-' tokens for testing
		protected.Use(middleware.AuthMiddleware(db))
		{
			// 0. Books
			books := protected.Group("/books")
			{
				books.POST("", bookHandler.CreateBook)
				books.GET("", bookHandler.GetBooks)
				books.POST("/:id/members", bookHandler.AddMember)
				books.GET("/:id/members", bookHandler.GetMembers)
			}

			// 1. Users
			users := protected.Group("/users")
			{
				users.GET("", userHandler.GetUsers)
				users.POST("", userHandler.CreateUser)
				users.GET("/me", userHandler.GetMe)
				users.PUT("/me", userHandler.UpdateMe)
				users.POST("/me/change-password", userHandler.ChangePassword)
				// Keep generic getters if needed by admin? Or just for compatibility
				users.GET("/:id", userHandler.GetUser)
			}

			// 2. Wallets
			wallets := protected.Group("/wallets")
			{
				wallets.GET("", walletHandler.GetWallets)
				wallets.POST("", walletHandler.CreateWallet)
				wallets.GET("/:id", walletHandler.GetWallet)
				wallets.PUT("/:id", walletHandler.UpdateWallet)
				wallets.DELETE("/:id", walletHandler.DeleteWallet)
			}

			// 3. Categories
			categories := protected.Group("/categories")
			{
				categories.GET("", categoryHandler.GetCategories)
				categories.POST("", categoryHandler.CreateCategory)
				categories.PUT("/:id", categoryHandler.UpdateCategory)
				categories.DELETE("/:id", categoryHandler.DeleteCategory)
			}

			// 4. Transactions
			transactions := protected.Group("/transactions")
			{
				transactions.GET("", transactionHandler.GetTransactions)
				transactions.POST("", transactionHandler.CreateTransaction)
				transactions.GET("/:id", transactionHandler.GetTransaction)
				transactions.PUT("/:id", transactionHandler.UpdateTransaction)
				transactions.DELETE("/:id", transactionHandler.DeleteTransaction)
			}

			// 5. Debts
			debts := protected.Group("/debts")
			{
				debts.GET("", debtHandler.GetDebts)
				debts.POST("", debtHandler.CreateDebt)
				debts.GET("/:id", debtHandler.GetDebt)
				debts.PUT("/:id", debtHandler.UpdateDebt)
				debts.DELETE("/:id", debtHandler.DeleteDebt)
				debts.POST("/:id/payment", debtHandler.MakePayment)
			}

			// 6. Reports
			reports := protected.Group("/reports")
			{
				reports.GET("/summary", reportHandler.GetSummary)
				reports.GET("/balance-history", reportHandler.GetBalanceHistory)
				reports.GET("/category-pie", reportHandler.GetCategoryPie)
				reports.GET("/daily-cashflow", reportHandler.GetDailyCashFlow)
			}

			// 7. System / Mock
			system := protected.Group("/upload")
			{
				system.POST("/image", systemHandler.UploadImage)
			}
			
			// History (Legacy/Extra)
			protected.GET("/history", historyHandler.GetHistoryLogs)
		}
	}
}
