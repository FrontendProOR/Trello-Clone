package middlewares

import (
	"net/http"
	"strings"
	"tasks-service/src/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a Gin middleware to protect routes with JWT authentication.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return
		}

		// The token is usually in the form of "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		// Validate the JWT token
		tokenString := parts[1]
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Attach the claims (user ID, email) to the request context
		c.Set("user_id", claims["user_id"])
		c.Set("email", claims["email"])

		// Continue processing the request
		c.Next()
	}
}
