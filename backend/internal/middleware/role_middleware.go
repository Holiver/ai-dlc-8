package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminMiddleware creates a middleware that requires admin role
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Check if user is admin
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Admin access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// EmployeeMiddleware creates a middleware that requires employee role
func EmployeeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Check if user is employee
		if role != "employee" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Employee access required",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleMiddleware creates a middleware that checks for specific roles
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user role from context (set by AuthMiddleware)
		role, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		userRole := role.(string)

		// Check if user's role is in allowed roles
		allowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				allowed = true
				break
			}
		}

		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
