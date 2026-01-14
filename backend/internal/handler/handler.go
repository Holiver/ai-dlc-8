package handler

import (
	"awsome-shop/internal/service"
)

// Handlers holds all handler instances
type Handlers struct {
	Auth              *AuthHandler
	User              *UserHandler
	Product           *ProductHandler
	Redemption        *RedemptionHandler
	Points            *PointsHandler
	AdminUser         *AdminUserHandler
	AdminProduct      *AdminProductHandler
	AdminPoints       *AdminPointsHandler
	AdminOrder        *AdminOrderHandler
	AdminReport       *AdminReportHandler
}

// NewHandlers creates and initializes all handlers
func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth:         NewAuthHandler(services.Auth),
		User:         NewUserHandler(services.User),
		Product:      NewProductHandler(services.Product),
		Redemption:   NewRedemptionHandler(services.Redemption),
		Points:       NewPointsHandler(services.Points),
		AdminUser:    NewAdminUserHandler(services.User),
		AdminProduct: NewAdminProductHandler(services.Product),
		AdminPoints:  NewAdminPointsHandler(services.Points),
		AdminOrder:   NewAdminOrderHandler(services.Redemption),
		AdminReport:  NewAdminReportHandler(services.Points, services.Redemption),
	}
}
