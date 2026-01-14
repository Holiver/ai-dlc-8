package models

import (
	"time"
)

// PointsTransaction represents a points transaction in the system
type PointsTransaction struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	UserID          uint             `gorm:"not null" json:"user_id"`
	User            User             `gorm:"foreignKey:UserID" json:"user,omitempty"`
	TransactionType string           `gorm:"type:enum('grant','deduct','redemption');not null" json:"transaction_type"`
	Amount          int              `gorm:"not null" json:"amount"`
	BalanceAfter    int              `gorm:"not null" json:"balance_after"`
	Reason          string           `gorm:"size:500" json:"reason"`
	OperatorID      *uint            `json:"operator_id"`
	Operator        *User            `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
	RelatedOrderID  *uint            `json:"related_order_id"`
	RelatedOrder    *RedemptionOrder `gorm:"foreignKey:RelatedOrderID" json:"related_order,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
}

// TableName specifies the table name for PointsTransaction model
func (PointsTransaction) TableName() string {
	return "points_transactions"
}
