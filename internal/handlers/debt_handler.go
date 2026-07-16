package handlers

import (
	"net/http"
	"time"

	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DebtHandler struct {
	db *gorm.DB
}

func NewDebtHandler(db *gorm.DB) *DebtHandler {
	return &DebtHandler{db: db}
}

// CreateDebt godoc
// @Summary Create a new debt record
// @Description Create a new debt (Lent or Borrowed)
// @Tags debts
// @Accept json
// @Produce json
// @Param debt body models.Debt true "Debt Data"
// @Success 201 {object} models.Debt
// @Router /debts [post]
func (h *DebtHandler) CreateDebt(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	var debt models.Debt
	if err := c.ShouldBindJSON(&debt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	debt.ID = uuid.New().String()
	debt.UserID = userID

	// Convert empty string pointer for BookID to nil
	if debt.BookID != nil && *debt.BookID == "" {
		debt.BookID = nil
	}

	if debt.BookID != nil && !canEditInBook(h.db, *debt.BookID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to add debts to this book"})
		return
	}

	if result := h.db.Create(&debt); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, debt)
}

// GetDebts godoc
// @Summary Get all debts
// @Description Get all debts (optional filter by type)
// @Tags debts
// @Produce json
// @Param type query string false "Debt Type (LENT/BORROWED)"
// @Param book_id query string false "Filter by Book ID"
// @Success 200 {array} models.Debt
// @Router /debts [get]
func (h *DebtHandler) GetDebts(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	debtType := c.Query("type")
	bookID := c.Query("book_id")
	var debts []models.Debt

	var query *gorm.DB
	if bookID != "" {
		if !isBookMember(h.db, bookID, userID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this book"})
			return
		}
		// Book members share visibility of all debts in the book, not just their own.
		query = h.db.Where("book_id = ?", bookID)
	} else {
		query = h.db.Where("user_id = ?", userID)
	}

	if debtType != "" {
		query = query.Where("type = ?", debtType)
	}

	if result := query.Find(&debts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, debts)
}

// GetDebt godoc
// @Summary Get a debt by ID
// @Description Get a debt by ID
// @Tags debts
// @Produce json
// @Param id path string true "Debt ID"
// @Success 200 {object} models.Debt
// @Router /debts/{id} [get]
func (h *DebtHandler) GetDebt(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	var debt models.Debt

	if result := h.db.First(&debt, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Debt not found"})
		return
	}

	isOwner := debt.UserID == userID
	isSharedMember := debt.BookID != nil && isBookMember(h.db, *debt.BookID, userID)
	if !isOwner && !isSharedMember {
		c.JSON(http.StatusNotFound, gin.H{"error": "Debt not found"})
		return
	}

	c.JSON(http.StatusOK, debt)
}

// UpdateDebt godoc
// @Summary Update a debt
// @Description Update a debt by ID
// @Tags debts
// @Accept json
// @Produce json
// @Param id path string true "Debt ID"
// @Param debt body models.Debt true "Debt Data"
// @Success 200 {object} models.Debt
// @Router /debts/{id} [put]
func (h *DebtHandler) UpdateDebt(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	var debt models.Debt

	if result := h.db.First(&debt, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Debt not found"})
		return
	}

	isOwner := debt.UserID == userID
	canEditViaBook := debt.BookID != nil && canEditInBook(h.db, *debt.BookID, userID)
	if !isOwner && !canEditViaBook {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to edit this debt"})
		return
	}

	originalUserID := debt.UserID
	if err := c.ShouldBindJSON(&debt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID and original owner are not changed
	debt.ID = id
	debt.UserID = originalUserID

	h.db.Save(&debt)
	c.JSON(http.StatusOK, debt)
}

// DeleteDebt godoc
// @Summary Delete a debt
// @Description Delete a debt by ID
// @Tags debts
// @Produce json
// @Param id path string true "Debt ID"
// @Success 200 {object} map[string]string
// @Router /debts/{id} [delete]
func (h *DebtHandler) DeleteDebt(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")

	var debt models.Debt
	if result := h.db.First(&debt, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Debt not found"})
		return
	}

	isOwner := debt.UserID == userID
	canEditViaBook := debt.BookID != nil && canEditInBook(h.db, *debt.BookID, userID)
	if !isOwner && !canEditViaBook {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this debt"})
		return
	}

	if result := h.db.Delete(&models.Debt{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Debt deleted successfully"})
}

type PaymentRequest struct {
	Amount float64 `json:"amount"`
}

// MakePayment godoc
// @Summary Make a payment for a debt
// @Description Deduct amount from debt
// @Tags debts
// @Accept json
// @Produce json
// @Param id path string true "Debt ID"
// @Param payment body PaymentRequest true "Payment Amount"
// @Success 200 {object} models.Debt
// @Router /debts/{id}/payment [post]
func (h *DebtHandler) MakePayment(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	var req PaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment amount must be positive"})
		return
	}

	tx := h.db.Begin()

	var debt models.Debt
	if result := tx.First(&debt, "id = ?", id); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Debt not found"})
		return
	}

	isOwner := debt.UserID == userID
	canEditViaBook := debt.BookID != nil && canEditInBook(h.db, *debt.BookID, userID)
	if !isOwner && !canEditViaBook {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to pay this debt"})
		return
	}

	debt.RemainingAmount -= req.Amount
	if debt.RemainingAmount < 0 {
		debt.RemainingAmount = 0
	}

	// Update installment plan if exists
	if debt.IsInstallment && debt.InstallmentPlan != nil {
		debt.InstallmentPlan.PaidMonths++
	}

	// Move the actual money: deduct from the linked wallet and record a transaction,
	// mirroring how CreateTransaction affects wallet balances.
	if debt.WalletID != nil {
		var wallet models.Wallet
		if err := tx.First(&wallet, "id = ?", *debt.WalletID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Linked wallet not found"})
			return
		}
		wallet.Balance -= req.Amount
		if err := tx.Save(&wallet).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update wallet balance"})
			return
		}

		description := "ชำระหนี้: " + debt.Title
		payment := models.Transaction{
			ID:          uuid.New().String(),
			WalletID:    *debt.WalletID,
			Amount:      req.Amount,
			Type:        models.TransactionTypeExpense,
			Description: description,
			Date:        time.Now(),
			CreatedByID: userID,
			BookID:      debt.BookID,
		}
		if err := tx.Create(&payment).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment transaction"})
			return
		}
	}

	if err := tx.Save(&debt).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update debt"})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, debt)
}
