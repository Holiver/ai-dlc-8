package models

import (
	"time"
)

// ProductPriceHistory represents product price change history
type ProductPriceHistory struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `gorm:"not null" json:"product_id"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	OldPoints  *int      `json:"old_points"`
	NewPoints  int       `gorm:"not null" json:"new_points"`
	OperatorID *uint     `json:"operator_id"`
	Operator   *User     `gorm:"foreignKey:OperatorID" json:"operator,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

// TableName specifies the table name for ProductPriceHistory model
func (ProductPriceHistory) TableName() string {
	return "product_price_history"
}
