package repository

import (
	"awsome-shop/internal/models"
	"errors"

	"gorm.io/gorm"
)

// RedemptionOrderRepository handles redemption order data access operations
type RedemptionOrderRepository struct {
	db *gorm.DB
}

// NewRedemptionOrderRepository creates a new RedemptionOrderRepository instance
func NewRedemptionOrderRepository(db *gorm.DB) *RedemptionOrderRepository {
	return &RedemptionOrderRepository{db: db}
}

// Create creates a new redemption order
func (r *RedemptionOrderRepository) Create(order *models.RedemptionOrder) error {
	return r.db.Create(order).Error
}

// GetByID retrieves a redemption order by ID
func (r *RedemptionOrderRepository) GetByID(id uint) (*models.RedemptionOrder, error) {
	var order models.RedemptionOrder
	err := r.db.Preload("User").Preload("Product").First(&order, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

// GetByOrderNumber retrieves a redemption order by order number
func (r *RedemptionOrderRepository) GetByOrderNumber(orderNumber string) (*models.RedemptionOrder, error) {
	var order models.RedemptionOrder
	err := r.db.Preload("User").Preload("Product").
		Where("order_number = ?", orderNumber).
		First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

// GetByUserID retrieves all orders for a specific user (ordered by created_at desc)
func (r *RedemptionOrderRepository) GetByUserID(userID uint) ([]models.RedemptionOrder, error) {
	var orders []models.RedemptionOrder
	err := r.db.Preload("Product").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// Update updates a redemption order
func (r *RedemptionOrderRepository) Update(order *models.RedemptionOrder) error {
	return r.db.Save(order).Error
}

// UpdateStatus updates an order's status
func (r *RedemptionOrderRepository) UpdateStatus(orderID uint, status string) error {
	if status != "preparing" && status != "delivered" {
		return errors.New("invalid status: must be 'preparing' or 'delivered'")
	}

	return r.db.Model(&models.RedemptionOrder{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

// BatchUpdateStatus updates multiple orders' status by order numbers
func (r *RedemptionOrderRepository) BatchUpdateStatus(orderNumbers []string, status string) error {
	if status != "preparing" && status != "delivered" {
		return errors.New("invalid status: must be 'preparing' or 'delivered'")
	}

	return r.db.Model(&models.RedemptionOrder{}).
		Where("order_number IN ?", orderNumbers).
		Update("status", status).Error
}

// List retrieves all orders with optional filters
func (r *RedemptionOrderRepository) List(status *string, userID *uint) ([]models.RedemptionOrder, error) {
	var orders []models.RedemptionOrder
	query := r.db.Preload("User").Preload("Product")

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	err := query.Order("created_at DESC").Find(&orders).Error
	return orders, err
}

// GetAll retrieves all orders (for admin)
func (r *RedemptionOrderRepository) GetAll() ([]models.RedemptionOrder, error) {
	var orders []models.RedemptionOrder
	err := r.db.Preload("User").Preload("Product").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// GetOrderStats retrieves order statistics for reports
type OrderStats struct {
	ProductName  string
	ProductID    uint
	UserName     string
	UserEmail    string
	PointsCost   int
	Status       string
	CreatedAt    string
}

func (r *RedemptionOrderRepository) GetOrderStats() ([]OrderStats, error) {
	var stats []OrderStats
	err := r.db.Table("redemption_orders").
		Select("redemption_orders.product_name, redemption_orders.product_id, users.full_name as user_name, users.email as user_email, redemption_orders.points_cost, redemption_orders.status, redemption_orders.created_at").
		Joins("LEFT JOIN users ON users.id = redemption_orders.user_id").
		Order("redemption_orders.created_at DESC").
		Scan(&stats).Error
	return stats, err
}

// CountByStatus counts orders by status
func (r *RedemptionOrderRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.RedemptionOrder{}).
		Where("status = ?", status).
		Count(&count).Error
	return count, err
}

// ExistsByOrderNumber checks if an order with the given order number exists
func (r *RedemptionOrderRepository) ExistsByOrderNumber(orderNumber string) (bool, error) {
	var count int64
	err := r.db.Model(&models.RedemptionOrder{}).
		Where("order_number = ?", orderNumber).
		Count(&count).Error
	return count > 0, err
}
