package handlers

import (
	"net/http"
	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionHandler struct {
	db *gorm.DB
}

func NewTransactionHandler(db *gorm.DB) *TransactionHandler {
	return &TransactionHandler{db: db}
}

// CreateTransaction godoc
// @Summary Create a new transaction
// @Description Create a new transaction (Income, Expense, Transfer)
// @Tags transactions
// @Accept json
// @Produce json
// @Param transaction body models.Transaction true "Transaction Data"
// @Success 201 {object} models.Transaction
// @Router /transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	var transaction models.Transaction
	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction.ID = uuid.New().String()

	// Set CreatedByID from context
	transaction.CreatedByID = c.MustGet("user_id").(string)

	// Convert empty string pointer to nil to avoid FK constraint violation
	if transaction.ToWalletID != nil && *transaction.ToWalletID == "" {
		transaction.ToWalletID = nil
	}
	
	// Convert empty string pointer for BookID to nil
	if transaction.BookID != nil && *transaction.BookID == "" {
		transaction.BookID = nil
	}

	// Start database transaction
	tx := h.db.Begin()

	if result := tx.Create(&transaction); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Update Wallet Balance
	var wallet models.Wallet
	if err := tx.First(&wallet, "id = ?", transaction.WalletID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wallet not found"})
		return
	}

	switch transaction.Type {
	case "INCOME":
		wallet.Balance += transaction.Amount
	case "EXPENSE":
		wallet.Balance -= transaction.Amount
	case "TRANSFER":
		wallet.Balance -= transaction.Amount
		
		if transaction.ToWalletID != nil {
			var toWallet models.Wallet
			if err := tx.First(&toWallet, "id = ?", *transaction.ToWalletID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Target wallet not found"})
				return
			}
			toWallet.Balance += transaction.Amount
			if err := tx.Save(&toWallet).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update target wallet balance"})
				return
			}
		}
	}

	if err := tx.Save(&wallet).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
		return
	}

	tx.Commit()

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactions godoc
// @Summary Get all transactions
// @Description Get all transactions with optional filters
// @Tags transactions
// @Produce json
// @Param wallet_id query string false "Filter by Wallet ID"
// @Param category_id query string false "Filter by Category ID"
// @Param start_date query string false "Start Date (YYYY-MM-DD)"
// @Param end_date query string false "End Date (YYYY-MM-DD)"
// @Param search query string false "Search Description"
// @Param book_id query string false "Filter by Book ID"
// @Success 200 {array} models.Transaction
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	walletID := c.Query("wallet_id")
	categoryID := c.Query("category_id")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	search := c.Query("search")
	bookID := c.Query("book_id")

	var transactions []models.Transaction

	query := h.db.Model(&models.Transaction{}).Where("created_by_id = ?", userID)
	
	if bookID != "" {
		query = query.Where("book_id = ?", bookID)
	}
	
	if walletID != "" {
		query = query.Where("wallet_id = ?", walletID)
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if startDate != "" && endDate != "" {
		query = query.Where("date BETWEEN ? AND ?", startDate, endDate)
	}
	if search != "" {
		query = query.Where("description LIKE ?", "%"+search+"%")
	}

	if result := query.Preload("Category").Preload("Wallet").Preload("ToWallet").Find(&transactions); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// GetTransaction godoc
// @Summary Get a transaction by ID
// @Description Get a transaction by ID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	var transaction models.Transaction

	if result := h.db.Preload("Category").Preload("Wallet").Preload("ToWallet").Where("id = ? AND created_by_id = ?", id, userID).First(&transaction); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateTransaction godoc
// @Summary Update a transaction
// @Description Update a transaction by ID
// @Tags transactions
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Param transaction body models.Transaction true "Transaction Data"
// @Success 200 {object} models.Transaction
// @Router /transactions/{id} [put]
func (h *TransactionHandler) UpdateTransaction(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	var transaction models.Transaction

	// Start DB Transaction
	tx := h.db.Begin()

	var existingTransaction models.Transaction
	if result := tx.Where("id = ? AND created_by_id = ?", id, userID).First(&existingTransaction); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// 1. Revert impact of existing transaction
	var oldWallet models.Wallet
	if err := tx.First(&oldWallet, "id = ?", existingTransaction.WalletID).Error; err == nil {
		switch existingTransaction.Type {
		case "INCOME":
			oldWallet.Balance -= existingTransaction.Amount
		case "EXPENSE":
			oldWallet.Balance += existingTransaction.Amount
		case "TRANSFER":
			oldWallet.Balance += existingTransaction.Amount
			if existingTransaction.ToWalletID != nil {
				var oldToWallet models.Wallet
				if err := tx.First(&oldToWallet, "id = ?", *existingTransaction.ToWalletID).Error; err == nil {
					oldToWallet.Balance -= existingTransaction.Amount
					tx.Save(&oldToWallet)
				}
			}
		}
		tx.Save(&oldWallet)
	}

	// 2. Bind new data (but keep IDs safe)
	if err := c.ShouldBindJSON(&transaction); err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction.ID = id
	transaction.CreatedByID = userID

	// Convert empty string pointer to nil for ToWalletID
	if transaction.ToWalletID != nil && *transaction.ToWalletID == "" {
		transaction.ToWalletID = nil
	}

	// 3. Apply impact of new transaction
	var newWallet models.Wallet
	if err := tx.First(&newWallet, "id = ?", transaction.WalletID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "New Wallet not found"})
		return
	}

	switch transaction.Type {
	case "INCOME":
		newWallet.Balance += transaction.Amount
	case "EXPENSE":
		newWallet.Balance -= transaction.Amount
	case "TRANSFER":
		newWallet.Balance -= transaction.Amount
		if transaction.ToWalletID != nil {
			var newToWallet models.Wallet
			if err := tx.First(&newToWallet, "id = ?", *transaction.ToWalletID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusInternalServerError, gin.H{"error": "New Target wallet not found"})
				return
			}
			newToWallet.Balance += transaction.Amount
			tx.Save(&newToWallet)
		}
	}
	tx.Save(&newWallet)

	// 4. Save Transaction
	if err := tx.Save(&transaction).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update transaction"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, transaction)
}

// DeleteTransaction godoc
// @Summary Delete a transaction
// @Description Delete a transaction by ID
// @Tags transactions
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} map[string]string
// @Router /transactions/{id} [delete]
func (h *TransactionHandler) DeleteTransaction(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	// Start DB Transaction
	tx := h.db.Begin()

	// Get transaction details first to know amount and wallets
	var transaction models.Transaction
	if result := tx.Where("id = ? AND created_by_id = ?", id, userID).First(&transaction); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Revert Wallet Balance
	var wallet models.Wallet
	if err := tx.First(&wallet, "id = ?", transaction.WalletID).Error; err == nil {
		switch transaction.Type {
		case "INCOME":
			wallet.Balance -= transaction.Amount
		case "EXPENSE":
			wallet.Balance += transaction.Amount
		case "TRANSFER":
			wallet.Balance += transaction.Amount
			
			if transaction.ToWalletID != nil {
				var toWallet models.Wallet
				if err := tx.First(&toWallet, "id = ?", *transaction.ToWalletID).Error; err == nil {
					toWallet.Balance -= transaction.Amount
					tx.Save(&toWallet)
				}
			}
		}
		tx.Save(&wallet)
	}

	// Delete Transaction
	if result := tx.Delete(&models.Transaction{}, "id = ?", id); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
