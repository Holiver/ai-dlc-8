package router

import (
	"awsome-shop/internal/config"
	"awsome-shop/internal/handler"
	"awsome-shop/internal/middleware"
	"awsome-shop/internal/repository"
	"awsome-shop/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Setup initializes the router with all routes
func Setup(db *gorm.DB, cfg *config.Config) *gin.Engine {
	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Create router without default middleware
	r := gin.New()

	// Add custom middleware
	r.Use(middleware.RecoveryMiddleware())
	r.Use(middleware.LoggingMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Initialize repositories
	repos := repository.NewRepositories(db)

	// Initialize services
	services := service.NewServices(repos, db, cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	// Initialize handlers
	handlers := handler.NewHandlers(services)

	// Create middleware instances
	authMiddleware := middleware.AuthMiddleware(services.Auth)
	adminMiddleware := middleware.AdminMiddleware()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Register all routes
		handlers.Auth.RegisterRoutes(v1, authMiddleware)
		handlers.User.RegisterRoutes(v1, authMiddleware)
		handlers.Product.RegisterRoutes(v1, authMiddleware)
		handlers.Redemption.RegisterRoutes(v1, authMiddleware)
		handlers.Points.RegisterRoutes(v1, authMiddleware)
		
		// Admin routes
		handlers.AdminUser.RegisterRoutes(v1, authMiddleware, adminMiddleware)
		handlers.AdminProduct.RegisterRoutes(v1, authMiddleware, adminMiddleware)
		handlers.AdminPoints.RegisterRoutes(v1, authMiddleware, adminMiddleware)
		handlers.AdminOrder.RegisterRoutes(v1, authMiddleware, adminMiddleware)
		handlers.AdminReport.RegisterRoutes(v1, authMiddleware, adminMiddleware)
	}

	return r
}
