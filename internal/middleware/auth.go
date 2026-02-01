package middleware

import (
	"raijai-backend/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AuthMiddleware - JWT Middleware
func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
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

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		// Optional: Check if user exists in DB and is active
		// var user models.User
		// if err := db.First(&user, "id = ?", claims.UserID).Error; err != nil {
		// 	c.AbortWithStatusJSON(401, gin.H{"error": "User not found or inactive"})
		// 	return
		// }

		c.Set("user_id", claims.UserID)
		c.Next()
	}
}
