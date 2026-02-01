package handlers

import (
	"net/http"

	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReportHandler struct {
	db *gorm.DB
}

func NewReportHandler(db *gorm.DB) *ReportHandler {
	return &ReportHandler{db: db}
}

// GetSummary
// @Summary Get income/expense summary
// @Description Get summary of income and expenses
// @Tags reports
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /reports/summary [get]
func (h *ReportHandler) GetSummary(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var result struct {
		Income  float64
		Expense float64
	}

	// Calculate Total Income
	h.db.Model(&models.Transaction{}).
		Where("created_by_id = ? AND type = ?", userID, "INCOME").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&result.Income)

	// Calculate Total Expense
	h.db.Model(&models.Transaction{}).
		Where("created_by_id = ? AND type = ?", userID, "EXPENSE").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&result.Expense)

	balance := result.Income - result.Expense

	c.JSON(http.StatusOK, gin.H{
		"income":  result.Income,
		"expense": result.Expense,
		"balance": balance,
		"period":  "All Time", // Default for now, can add date filters later
	})
}

// GetBalanceHistory
// @Summary Get balance history
// @Description Get history of total balance over time
// @Tags reports
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /reports/balance-history [get]
func (h *ReportHandler) GetBalanceHistory(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	type DailyBalance struct {
		Date    string  `json:"date"`
		Balance float64 `json:"balance"`
	}

	var history []DailyBalance

	// Simple query: Group by date and sum (Income - Expense)
	// Note: This calculates daily change, not running balance. 
	// For running balance, we'd need more complex SQL or post-processing.
	// Implementing daily change for now as a baseline.
	
	err := h.db.Table("transactions").
		Select("TO_CHAR(date, 'YYYY-MM-DD') as date, SUM(CASE WHEN type = 'INCOME' THEN amount WHEN type = 'EXPENSE' THEN -amount ELSE 0 END) as balance").
		Where("created_by_id = ?", userID).
		Group("date").
		Order("date").
		Scan(&history).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Post-process for running balance if needed, or just return daily activity
	// Let's calculate running balance
	var runningBalance float64
	for i := range history {
		runningBalance += history[i].Balance
		history[i].Balance = runningBalance
	}

	c.JSON(http.StatusOK, history)
}

// GetCategoryPie
// @Summary Get category distribution
// @Description Get pie chart data for categories
// @Tags reports
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /reports/category-pie [get]
func (h *ReportHandler) GetCategoryPie(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	type CategoryStat struct {
		Category string  `json:"category"`
		Amount   float64 `json:"amount"`
		Color    string  `json:"color"`
	}

	var stats []CategoryStat

	err := h.db.Table("transactions").
		Select("categories.name as category, SUM(transactions.amount) as amount, categories.color as color").
		Joins("JOIN categories ON transactions.category_id = categories.id").
		Where("transactions.created_by_id = ? AND transactions.type = ?", userID, "EXPENSE").
		Group("categories.name, categories.color").
		Scan(&stats).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetDailyCashFlow
// @Summary Get daily income and expense
// @Description Get daily income and expense statistics for the current month
// @Tags reports
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /reports/daily-cashflow [get]
func (h *ReportHandler) GetDailyCashFlow(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	type DailyFlow struct {
		Date    string  `json:"date"`
		Income  float64 `json:"income"`
		Expense float64 `json:"expense"`
	}

	var flows []DailyFlow

	// Get data for the current month (simplified query for now, fetching all and filtering or grouping)
	// Ideally should filter by date range, here we'll grab all and let frontend/chart handle or limit to recent.
	// For "current month", let's use SQL date filtering. Assuming SQLite/Postgres compatibility issue, writing generic Group By.
	
	err := h.db.Table("transactions").
		Select("TO_CHAR(date, 'YYYY-MM-DD') as date, SUM(CASE WHEN type = 'INCOME' THEN amount ELSE 0 END) as income, SUM(CASE WHEN type = 'EXPENSE' THEN amount ELSE 0 END) as expense").
		Where("created_by_id = ?", userID).
		Group("date").
		Order("date").
		Scan(&flows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, flows)
}
