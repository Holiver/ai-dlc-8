package service

import (
	"awsome-shop/internal/repository"

	"gorm.io/gorm"
)

// Services holds all service instances
type Services struct {
	Auth       *AuthService
	User       *UserService
	Product    *ProductService
	Points     *PointsService
	Redemption *RedemptionService
}

// NewServices creates and initializes all services
func NewServices(repos *repository.Repositories, db *gorm.DB, jwtSecret string, jwtExpirationHours int) *Services {
	// Create AuthService first (needed by UserService)
	authService := NewAuthService(
		repos.User,
		repos.PointsTransaction,
		db,
		jwtSecret,
		jwtExpirationHours,
	)

	// Create other services
	userService := NewUserService(
		repos.User,
		repos.PointsTransaction,
		authService,
		db,
	)

	productService := NewProductService(
		repos.Product,
		db,
	)

	pointsService := NewPointsService(
		repos.User,
		repos.PointsTransaction,
		db,
	)

	redemptionService := NewRedemptionService(
		repos.User,
		repos.Product,
		repos.RedemptionOrder,
		repos.PointsTransaction,
		db,
	)

	return &Services{
		Auth:       authService,
		User:       userService,
		Product:    productService,
		Points:     pointsService,
		Redemption: redemptionService,
	}
}
