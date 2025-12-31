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
	if userID, exists := c.Get("user_id"); exists {
		transaction.CreatedByID = userID.(string)
	}

	// Convert empty string pointer to nil to avoid FK constraint violation
	if transaction.ToWalletID != nil && *transaction.ToWalletID == "" {
		transaction.ToWalletID = nil
	}

	if result := h.db.Create(&transaction); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// GetTransactions godoc
// @Summary Get all transactions
// @Description Get all transactions
// @Tags transactions
// @Produce json
// @Param wallet_id query string false "Filter by Wallet ID"
// @Success 200 {array} models.Transaction
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	walletID := c.Query("wallet_id")
	var transactions []models.Transaction

	query := h.db.Model(&models.Transaction{})
	if walletID != "" {
		query = query.Where("wallet_id = ?", walletID)
	}

	if result := query.Find(&transactions); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, transactions)
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
	id := c.Param("id")
	var transaction models.Transaction

	if result := h.db.First(&transaction, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.db.Save(&transaction)
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
	id := c.Param("id")
	if result := h.db.Delete(&models.Transaction{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
