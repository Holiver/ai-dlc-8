package models

import (
	"time"
)

// RedemptionOrder represents a redemption order in the system
type RedemptionOrder struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	OrderNumber        string    `gorm:"uniqueIndex;size:50;not null" json:"order_number"`
	UserID             uint      `gorm:"not null" json:"user_id"`
	User               User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	ProductID          uint      `gorm:"not null" json:"product_id"`
	Product            Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	ProductName        string    `gorm:"size:200;not null" json:"product_name"`
	PointsCost         int       `gorm:"not null" json:"points_cost"`
	PointsBalanceAfter int       `gorm:"not null" json:"points_balance_after"`
	Status             string    `gorm:"type:enum('preparing','delivered');default:'preparing'" json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// TableName specifies the table name for RedemptionOrder model
func (RedemptionOrder) TableName() string {
	return "redemption_orders"
}
