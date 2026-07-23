package handlers

import (
	"net/http"
	"strings"
	"time"

	"raijai-backend/internal/models"
	"raijai-backend/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	var credentials struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := h.db.Where("username = ?", credentials.Username).First(&user); result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	if !utils.CheckPasswordHash(credentials.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// In a real app, refresh token logic would be more complex
	refreshToken := "mock-refresh-token-" + uuid.New().String()

	c.JSON(http.StatusOK, gin.H{
		"token":        token,
		"refreshToken": refreshToken,
		"expiresIn":    86400, // 24 hours
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
	var req models.UserRegistration
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if username already exists
	var existingUser models.User
	if result := h.db.Where("username = ?", req.Username).First(&existingUser); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUser := models.User{
		ID:        uuid.New().String(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Name:      req.FirstName + " " + req.LastName,
		AvatarURL: "https://ui-avatars.com/api/?name=" + req.Username,
		CreatedAt: time.Now(),
	}

	if result := h.db.Create(&newUser); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

// RefreshToken
// @Summary Refresh access token
// @Description Get a new access token using the current Bearer token (must still be valid/unexpired)
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "New Token"
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
		return
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return
	}

	claims, err := utils.ValidateToken(parts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Verify the user still exists
	var user models.User
	if err := h.db.First(&user, "id = ?", claims.UserID).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	newToken, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	refreshToken := "mock-refresh-token-" + uuid.New().String()

	c.JSON(http.StatusOK, gin.H{
		"token":        newToken,
		"refreshToken": refreshToken,
		"expiresIn":    86400,
	})
}

// ForgotPassword
// @Summary Request a password reset token
// @Description Generates a single-use password reset token for the given email.
// @Description DEV MODE: no email is sent — the raw token is returned directly
// @Description in the JSON response so the frontend can build the reset link.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body map[string]string true "Email"
// @Success 200 {object} map[string]string "message, and resetToken when the account exists"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if result := h.db.Where("email = ?", req.Email).First(&user); result.Error != nil {
		// TODO: once real email sending is added, this branch and the
		// success branch below must return an identical response (no
		// resetToken ever in the body) so account existence isn't leaked.
		c.JSON(http.StatusOK, gin.H{
			"message": "If this email is registered, a password reset token has been generated.",
		})
		return
	}

	h.db.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used = ?", user.ID, false).
		Update("used", true)

	rawToken, err := utils.GenerateResetToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate reset token"})
		return
	}

	resetToken := models.PasswordResetToken{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		TokenHash: utils.HashToken(rawToken),
		ExpiresAt: time.Now().Add(30 * time.Minute),
		CreatedAt: time.Now(),
	}

	if result := h.db.Create(&resetToken); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reset token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Password reset token generated (dev mode – no email is sent).",
		"resetToken": rawToken,
	})
}

// ResetPassword
// @Summary Reset password using a reset token
// @Description Consumes a password reset token (returned by /auth/forgot-password
// @Description in dev mode) to set a new password. Single-use; expires after 30 minutes.
// @Tags auth
// @Accept json
// @Produce json
// @Param body body map[string]string true "Token and NewPassword"
// @Success 200 {object} map[string]string
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req struct {
		Token       string `json:"token" binding:"required"`
		NewPassword string `json:"newPassword" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenHash := utils.HashToken(req.Token)

	var resetToken models.PasswordResetToken
	if result := h.db.Where("token_hash = ? AND used = ?", tokenHash, false).First(&resetToken); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	if time.Now().After(resetToken.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired reset token"})
		return
	}

	var user models.User
	if result := h.db.First(&user, "id = ?", resetToken.UserID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	user.Password = hashedPassword
	if result := h.db.Save(&user); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	h.db.Model(&models.PasswordResetToken{}).
		Where("user_id = ? AND used = ?", user.ID, false).
		Update("used", true)

	c.JSON(http.StatusOK, gin.H{"message": "Password has been reset successfully"})
}
