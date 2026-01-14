package models

import (
	"time"
)

// Product represents a product in the system
type Product struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:200;not null" json:"name"`
	ImageURL       string    `gorm:"size:500" json:"image_url"`
	PointsRequired int       `gorm:"not null" json:"points_required"`
	StockQuantity  int       `gorm:"default:0" json:"stock_quantity"`
	Status         string    `gorm:"type:enum('active','inactive');default:'active'" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}
