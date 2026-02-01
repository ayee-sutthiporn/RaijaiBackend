package handlers

import (
	"net/http"
	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true "User Data"
// @Success 201 {object} models.User
// @Router /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.db.Create(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Get user details
// @Tags users
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if result := h.db.First(&user, "id = ?", id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetMe
// @Summary Get current user profile
// @Description Get profile of the logged-in user
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if result := h.db.First(&user, "id = ?", userID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

// UpdateMe
// @Summary Update current user profile
// @Description Update profile details of the logged-in user
// @Tags users
// @Accept json
// @Produce json
// @Param user body map[string]string true "User Details"
// @Success 200 {object} models.User
// @Router /users/me [put]
func (h *UserHandler) UpdateMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if result := h.db.First(&user, "id = ?", userID); result.Error != nil {
		// Mock Data if not found
		user = models.User{
			ID: userID.(string),
			Name: "Updated Name Mock",
			Email: "mock@example.com",
		}
	} else {
		user.Name = "Updated Name Mock"
	}
	
	c.JSON(http.StatusOK, user)
}

// ChangePassword
// @Summary Change password
// @Description Change the password of the logged-in user
// @Tags users
// @Accept json
// @Produce json
// @Param body body map[string]string true "Old and New Password"
// @Success 200 {object} map[string]string
// @Router /users/me/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Password changed successfully (Mock)",
	})
}
