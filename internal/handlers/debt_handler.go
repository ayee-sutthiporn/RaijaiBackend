package handlers

import (
	"net/http"
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

	if result := h.db.Create(&debt); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, debt)
}

// GetDebts godoc
// @Summary Get all debts
// @Description Get all debts
// @Tags debts
// @Produce json
// @Success 200 {array} models.Debt
// @Router /debts [get]
func (h *DebtHandler) GetDebts(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	var debts []models.Debt
	if result := h.db.Where("user_id = ?", userID).Find(&debts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, debts)
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

	if result := h.db.Where("id = ? AND user_id = ?", id, userID).First(&debt); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Debt not found"})
		return
	}

	if err := c.ShouldBindJSON(&debt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID and UserID are not changed
	debt.ID = id
	debt.UserID = userID

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
	if result := h.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Debt{}); result.Error != nil {
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

	var debt models.Debt
	if result := h.db.Where("id = ? AND user_id = ?", id, userID).First(&debt); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Debt not found"})
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

	h.db.Save(&debt)
	c.JSON(http.StatusOK, debt)
}
