package middlewares

import (
	"strings"

	"task-5-pbi-btpns-daffajatmiko/app/auth"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a function to protect routes by validating JWT tokens.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") // Get bearer token from request header
		if tokenString == "" {
			c.JSON(401, gin.H{"error": "Token not found"}) // Respond with an error if token is missing
			c.Abort()
			return
		}

		err := auth.ValidateToken(strings.Split(tokenString, "Bearer ")[1]) // Validate the token
		if err != nil {
			c.JSON(401, gin.H{"error": err.Error()}) // Respond with an error if token validation fails
			c.Abort()
			return
		}
		c.Next() // Continue to the next middleware or handler if token is valid
	}
}
