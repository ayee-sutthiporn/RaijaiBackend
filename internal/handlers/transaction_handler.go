package handlers

import (
	"net/http"
	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
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
