package handlers

import (
	"net/http"
	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CategoryHandler struct {
	db *gorm.DB
}

func NewCategoryHandler(db *gorm.DB) *CategoryHandler {
	return &CategoryHandler{db: db}
}

// CreateCategory godoc
// @Summary Create a new category
// @Description Create a new category
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category Data"
// @Success 201 {object} models.Category
// @Router /categories [post]
func (h *CategoryHandler) CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.db.Create(&category); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetCategories godoc
// @Summary Get all categories
// @Description Get all categories
// @Tags categories
// @Produce json
// @Success 200 {array} models.Category
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	var categories []models.Category
	if result := h.db.Find(&categories); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}
