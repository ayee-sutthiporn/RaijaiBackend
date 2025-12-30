package middleware

import (
	"context"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(issuerURL, clientID string) gin.HandlerFunc {
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

		provider, err := oidc.NewProvider(context.Background(), issuerURL)
		if err != nil {
			// In development, we might not reach Keycloak. Log error but maybe fail safe?
			// For now, fail hard.
			c.AbortWithStatusJSON(500, gin.H{"error": "Failed to connect to identity provider: " + err.Error()})
			return
		}

		verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
		idToken, err := verifier.Verify(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token: " + err.Error()})
			return
		}

		c.Set("user_id", idToken.Subject)

		// Token is valid
		c.Next()
	}
}
