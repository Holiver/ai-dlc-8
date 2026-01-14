package repository

import (
	"gorm.io/gorm"
)

// Repositories holds all repository instances
type Repositories struct {
	User              *UserRepository
	Product           *ProductRepository
	RedemptionOrder   *RedemptionOrderRepository
	PointsTransaction *PointsTransactionRepository
}

// NewRepositories creates and initializes all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:              NewUserRepository(db),
		Product:           NewProductRepository(db),
		RedemptionOrder:   NewRedemptionOrderRepository(db),
		PointsTransaction: NewPointsTransactionRepository(db),
	}
}
