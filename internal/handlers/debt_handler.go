package handlers

import (
	"net/http"
	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
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
	var debt models.Debt
	if err := c.ShouldBindJSON(&debt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
// @Description Get all debts
// @Tags debts
// @Produce json
// @Success 200 {array} models.Debt
// @Router /debts [get]
func (h *DebtHandler) GetDebts(c *gin.Context) {
	var debts []models.Debt
	if result := h.db.Find(&debts); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, debts)
}
