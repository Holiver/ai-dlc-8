package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware creates a middleware for panic recovery
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error
				log.Printf("Panic recovered: %v", err)

				// Return error response
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})

				c.Abort()
			}
		}()

		c.Next()
	}
}
