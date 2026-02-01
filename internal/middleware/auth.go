package middleware

import (
	"context"
	"raijai-backend/internal/models"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// AuthMiddleware - Original OIDC Middleware (kept for reference, but we might switch to Mock)
func AuthMiddleware(issuerURL, clientID string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header required"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenString := parts[1]

		// --- MOCK OVERRIDE START ---
		// If token starts with "mock-", treat as mock authenticated user
		if strings.HasPrefix(tokenString, "mock-") {
			userID := "mock-user-id-123"
			// Ensure mock user exists in DB
			db.FirstOrCreate(&models.User{
				ID:    userID,
				Name:  "Mock User",
				Email: "mock@example.com",
			})
			c.Set("user_id", userID)
			c.Next()
			return
		}
		// --- MOCK OVERRIDE END ---

		provider, err := oidc.NewProvider(context.Background(), issuerURL)
		if err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Failed to connect to identity provider: " + err.Error()})
			return
		}

		verifier := provider.Verifier(&oidc.Config{SkipClientIDCheck: true})
		idToken, err := verifier.Verify(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		// Sync user to local DB
		var claims struct {
			Sub   string `json:"sub"`
			Email string `json:"email"`
			Name  string `json:"name"`
		}
		if err := idToken.Claims(&claims); err != nil {
			c.AbortWithStatusJSON(500, gin.H{"error": "Failed to parse claims"})
			return
		}

		user := models.User{
			ID:    claims.Sub,
			Email: claims.Email,
			Name:  claims.Name,
		}

		// Upsert user
		if result := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"email", "name"}),
		}).Create(&user); result.Error != nil {
			// Fallback if upsert fails
			var existingUser models.User
			if err := db.FirstOrCreate(&existingUser, models.User{ID: user.ID}).Error; err != nil {
				c.AbortWithStatusJSON(500, gin.H{"error": "Failed to sync user: " + err.Error()})
				return
			}
			db.Model(&existingUser).Updates(models.User{Email: user.Email, Name: user.Name})
		}

		c.Set("user_id", idToken.Subject)
		c.Next()
	}
}
