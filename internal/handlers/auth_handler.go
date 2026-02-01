package handlers

import (
	"net/http"
	"time"

	"raijai-backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthHandler struct {
	db *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Login
// @Summary Login to the system
// @Description Login with username and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body map[string]string true "Username and Password"
// @Success 200 {object} map[string]string "Token"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	// Mock User Data
	user := models.User{
		ID:        "1",
		Username:  "demo_user",
		Email:     "demo@raijai.com",
		FirstName: "Demo",
		LastName:  "User",
		Name:      "Demo User",
		AvatarURL: "https://ui-avatars.com/api/?name=Demo+User",
		CreatedAt: time.Now(),
	}

	// Mock Login Response
	c.JSON(http.StatusOK, gin.H{
		"token":        "mock-jwt-token-example-123456",
		"refreshToken": "mock-refresh-token-example-123456",
		"expiresIn":    3600,
		"user":         user,
	})
}

// Register
// @Summary Register a new user
// @Description Register a new user account
// @Tags auth
// @Accept json
// @Produce json
// @Param user body map[string]string true "User Details"
// @Success 201 {object} map[string]string "Success Message"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// Mock Register
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully (Mock)",
	})
}

// RefreshToken
// @Summary Refresh access token
// @Description Get a new access token using a refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param body body map[string]string true "Refresh Token"
// @Success 200 {object} map[string]string "New Token"
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Mock Refresh Token
	c.JSON(http.StatusOK, gin.H{
		"token":     "new-mock-jwt-token-" + time.Now().Format("150405"),
		"expiresIn": 3600,
	})
}
