package router

import (
	"awsome-shop/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup initializes the router with all routes
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "login endpoint"})
			})
			auth.POST("/logout", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "logout endpoint"})
			})
			auth.GET("/me", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "me endpoint"})
			})
		}

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/profile", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "profile endpoint"})
			})
			users.PUT("/phone", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "update phone endpoint"})
			})
		}

		// Product routes
		products := v1.Group("/products")
		{
			products.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "products list endpoint"})
			})
		}

		// Redemption routes
		redemptions := v1.Group("/redemptions")
		{
			redemptions.POST("", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "create redemption endpoint"})
			})
			redemptions.GET("", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "redemption history endpoint"})
			})
		}

		// Points routes
		points := v1.Group("/points")
		{
			points.GET("/balance", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "points balance endpoint"})
			})
			points.GET("/transactions", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "points transactions endpoint"})
			})
		}

		// Admin routes
		admin := v1.Group("/admin")
		{
			// Admin user management
			adminUsers := admin.Group("/users")
			{
				adminUsers.POST("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "create user endpoint"})
				})
				adminUsers.PUT("/:id/status", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "update user status endpoint"})
				})
			}

			// Admin product management
			adminProducts := admin.Group("/products")
			{
				adminProducts.POST("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "create product endpoint"})
				})
				adminProducts.PUT("/:id", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "update product endpoint"})
				})
				adminProducts.PUT("/:id/status", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "update product status endpoint"})
				})
				adminProducts.POST("/batch", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "batch import products endpoint"})
				})
			}

			// Admin points management
			adminPoints := admin.Group("/points")
			{
				adminPoints.POST("/grant", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "grant points endpoint"})
				})
				adminPoints.POST("/deduct", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "deduct points endpoint"})
				})
				adminPoints.POST("/batch-grant", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "batch grant points endpoint"})
				})
			}

			// Admin order management
			adminOrders := admin.Group("/orders")
			{
				adminOrders.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "list orders endpoint"})
				})
				adminOrders.PUT("/batch-status", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "batch update order status endpoint"})
				})
			}

			// Admin reports
			adminReports := admin.Group("/reports")
			{
				adminReports.GET("/points-grants", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "points grants report endpoint"})
				})
				adminReports.GET("/points-balances", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "points balances report endpoint"})
				})
				adminReports.GET("/redemptions", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "redemptions report endpoint"})
				})
			}
		}
	}

	return r
}
