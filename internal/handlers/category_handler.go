package handlers

import (
	"net/http"
	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	userID := c.MustGet("user_id").(string)
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	category.ID = uuid.New().String()
	category.UserID = userID

	// Convert empty string pointer for BookID to nil
	if category.BookID != nil && *category.BookID == "" {
		category.BookID = nil
	}

	if category.BookID != nil && !canEditInBook(h.db, *category.BookID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to add categories to this book"})
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
// @Description Get all categories (optional filter by type)
// @Tags categories
// @Produce json
// @Param type query string false "Category Type (INCOME/EXPENSE)"
// @Param book_id query string false "Filter by Book ID"
// @Success 200 {array} models.Category
// @Router /categories [get]
func (h *CategoryHandler) GetCategories(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	categoryType := c.Query("type")
	bookID := c.Query("book_id")
	var categories []models.Category

	var query *gorm.DB
	if bookID != "" {
		if !isBookMember(h.db, bookID, userID) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not a member of this book"})
			return
		}
		// Book members share visibility of all categories in the book, not just their own.
		query = h.db.Where("book_id = ?", bookID)
	} else {
		query = h.db.Where("user_id = ?", userID)
	}

	if categoryType != "" {
		query = query.Where("type = ?", categoryType)
	}

	if result := query.Find(&categories); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update a category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.Category true "Category Data"
// @Success 200 {object} models.Category
// @Router /categories/{id} [put]
func (h *CategoryHandler) UpdateCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")
	var category models.Category

	if result := h.db.First(&category, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	isOwner := category.UserID == userID
	canEditViaBook := category.BookID != nil && canEditInBook(h.db, *category.BookID, userID)
	if !isOwner && !canEditViaBook {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to edit this category"})
		return
	}

	originalUserID := category.UserID
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure ID and original owner are not changed
	category.ID = id
	category.UserID = originalUserID

	h.db.Save(&category)
	c.JSON(http.StatusOK, category)
}

// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete a category by ID
// @Tags categories
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} map[string]string
// @Router /categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	id := c.Param("id")

	var category models.Category
	if result := h.db.First(&category, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	isOwner := category.UserID == userID
	canEditViaBook := category.BookID != nil && canEditInBook(h.db, *category.BookID, userID)
	if !isOwner && !canEditViaBook {
		c.JSON(http.StatusForbidden, gin.H{"error": "You do not have permission to delete this category"})
		return
	}

	if result := h.db.Delete(&models.Category{}, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}
