package handlers

import (
	"net/http"
	"time"

	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BookHandler struct {
	db *gorm.DB
}

func NewBookHandler(db *gorm.DB) *BookHandler {
	return &BookHandler{db: db}
}

// CreateBook
// @Summary Create a new book
// @Description Create a new ledger book
// @Tags books
// @Accept json
// @Produce json
// @Success 201 {object} models.Book
// @Router /books [post]
func (h *BookHandler) CreateBook(c *gin.Context) {
	userId := c.MustGet("user_id").(string)

	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	book := models.Book{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Transaction to create book and add owner as member
	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&book).Error; err != nil {
			return err
		}
		
		member := models.BookMember{
			BookID:   book.ID,
			UserID:   userId,
			Role:     "OWNER",
			JoinedAt: time.Now(),
		}
		if err := tx.Create(&member).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, book)
}

// GetBooks
// @Summary Get user's books
// @Description List all books the user is a member of
// @Tags books
// @Produce json
// @Success 200 {array} models.Book
// @Router /books [get]
func (h *BookHandler) GetBooks(c *gin.Context) {
	userId := c.MustGet("user_id").(string)

	var books []models.Book
	err := h.db.Joins("JOIN book_members on book_members.book_id = books.id").
		Where("book_members.user_id = ?", userId).
		Find(&books).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, books)
}

// AddMember
// @Summary Add member to book
// @Description Add a user to a book using email
// @Tags books
// @Accept json
// @Success 200 {object} map[string]string
// @Router /books/{id}/members [post]
func (h *BookHandler) AddMember(c *gin.Context) {
	bookID := c.Param("id")
	// TODO: verify requester permissions

	var req struct {
		Email string `json:"email" binding:"required,email"`
		Role  string `json:"role"` // EDITOR, VIEWER
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find User
	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Check if already member
	var existingMember models.BookMember
	if err := h.db.Where("book_id = ? AND user_id = ?", bookID, user.ID).First(&existingMember).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User is already a member"})
		return
	}

	member := models.BookMember{
		BookID:   bookID,
		UserID:   user.ID,
		Role:     req.Role, // Default to VIEWER if empty?
		JoinedAt: time.Now(),
	}
	if member.Role == "" { member.Role = "EDITOR" }

	if err := h.db.Create(&member).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

// GetMembers
// @Summary Get book members
// @Description List all members of a book
// @Tags books
// @Produce json
// @Success 200 {array} models.User
// @Router /books/{id}/members [get]
func (h *BookHandler) GetMembers(c *gin.Context) {
	bookID := c.Param("id")

	// Join User and BookMember to get Role + User details
	// This is a bit tricky with GORM structs vs flat json
	// Simplied: Just fetch users for now
	var users []models.User
	err := h.db.Table("users").
		Joins("JOIN book_members on book_members.user_id = users.id").
		Where("book_members.book_id = ?", bookID).
		Scan(&users).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// DeleteBook
// @Summary Delete a book
// @Description Delete a book by ID (Only owner can delete)
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} map[string]string
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(c *gin.Context) {
	userId := c.MustGet("user_id").(string)
	bookId := c.Param("id")

	// Check if user is owner
	var book models.Book
	if err := h.db.Where("id = ? AND owner_id = ?", bookId, userId).First(&book).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "Book not found or you are not the owner"})
		return
	}

	// Delete Book
	if err := h.db.Delete(&book).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}
