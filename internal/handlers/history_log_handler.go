package handlers

import (
	"net/http"
	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HistoryLogHandler struct {
	db *gorm.DB
}

func NewHistoryLogHandler(db *gorm.DB) *HistoryLogHandler {
	return &HistoryLogHandler{db: db}
}

// GetHistoryLogs godoc
// @Summary Get history logs
// @Description Get history logs, optionally filtered by entity
// @Tags history
// @Produce json
// @Param entity_id query string false "Filter by Entity ID"
// @Param entity_type query string false "Filter by Entity Type"
// @Success 200 {array} models.HistoryLog
// @Router /history [get]
func (h *HistoryLogHandler) GetHistoryLogs(c *gin.Context) {
	entityID := c.Query("entity_id")
	entityType := c.Query("entity_type")
	var logs []models.HistoryLog

	query := h.db.Model(&models.HistoryLog{})
	if entityID != "" {
		query = query.Where("entity_id = ?", entityID)
	}
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}

	if result := query.Find(&logs); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, logs)
}
