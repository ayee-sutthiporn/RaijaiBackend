package handlers

import (
	"net/http"

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
	// Mock Summary Data
	c.JSON(http.StatusOK, gin.H{
		"income":  25000.00,
		"expense": 12500.50,
		"balance": 12499.50,
		"period":  "This Month",
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
	// Mock History Data
	c.JSON(http.StatusOK, []gin.H{
		{"date": "2023-10-01", "balance": 10000},
		{"date": "2023-10-05", "balance": 15000},
		{"date": "2023-10-10", "balance": 12000},
		{"date": "2023-10-15", "balance": 20000},
		{"date": "2023-10-20", "balance": 18000},
	})
}

// GetCategoryPie
// @Summary Get category distribution
// @Description Get pie chart data for categories
// @Tags reports
// @Produce json
// @Success 200 {array} map[string]interface{}
// @Router /reports/category-pie [get]
func (h *ReportHandler) GetCategoryPie(c *gin.Context) {
	// Mock Pie Chart Data
	c.JSON(http.StatusOK, []gin.H{
		{"category": "Food", "amount": 5000, "color": "#FF6384"},
		{"category": "Transport", "amount": 3000, "color": "#36A2EB"},
		{"category": "Entertainment", "amount": 2000, "color": "#FFCE56"},
		{"category": "Utilities", "amount": 2500, "color": "#4BC0C0"},
	})
}
