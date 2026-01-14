package service

import (
	"awsome-shop/internal/models"
	"awsome-shop/internal/repository"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

// RedemptionService handles product redemption operations
type RedemptionService struct {
	userRepo              *repository.UserRepository
	productRepo           *repository.ProductRepository
	orderRepo             *repository.RedemptionOrderRepository
	pointsTransactionRepo *repository.PointsTransactionRepository
	db                    *gorm.DB
}

// NewRedemptionService creates a new RedemptionService instance
func NewRedemptionService(
	userRepo *repository.UserRepository,
	productRepo *repository.ProductRepository,
	orderRepo *repository.RedemptionOrderRepository,
	pointsTransactionRepo *repository.PointsTransactionRepository,
	db *gorm.DB,
) *RedemptionService {
	return &RedemptionService{
		userRepo:              userRepo,
		productRepo:           productRepo,
		orderRepo:             orderRepo,
		pointsTransactionRepo: pointsTransactionRepo,
		db:                    db,
	}
}

// RedeemProductRequest represents a request to redeem a product
type RedeemProductRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

// RedeemProduct processes a product redemption
// This includes: points validation, stock validation, points deduction, stock reduction, order creation
func (s *RedemptionService) RedeemProduct(userID uint, productID uint) (*models.RedemptionOrder, error) {
	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user with lock
	var user models.User
	err := tx.Clauses(gorm.Locking{Strength: "UPDATE"}).First(&user, userID).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("user not found")
	}

	// Check if user is active
	if !user.IsActive {
		tx.Rollback()
		return nil, errors.New("user account is inactive")
	}

	// Get product with lock (for stock management)
	var product models.Product
	err = tx.Clauses(gorm.Locking{Strength: "UPDATE"}).First(&product, productID).Error
	if err != nil {
		tx.Rollback()
		return nil, errors.New("product not found")
	}

	// Validate product is active
	if product.Status != "active" {
		tx.Rollback()
		return nil, errors.New("product is not available")
	}

	// Validate stock
	if product.StockQuantity <= 0 {
		tx.Rollback()
		return nil, errors.New("product is out of stock")
	}

	// Validate user has sufficient points
	if user.PointsBalance < product.PointsRequired {
		tx.Rollback()
		return nil, errors.New("insufficient points")
	}

	// Calculate new balance
	newBalance := user.PointsBalance - product.PointsRequired

	// Update user points balance
	err = tx.Model(&models.User{}).
		Where("id = ?", userID).
		Update("points_balance", newBalance).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Decrement product stock
	err = tx.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", 1)).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Generate order number
	orderNumber := s.generateOrderNumber(userID, productID)

	// Create redemption order
	order := &models.RedemptionOrder{
		OrderNumber:        orderNumber,
		UserID:             userID,
		ProductID:          productID,
		ProductName:        product.Name,
		PointsCost:         product.PointsRequired,
		PointsBalanceAfter: newBalance,
		Status:             "preparing",
	}

	err = tx.Create(order).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create points transaction record
	transaction := &models.PointsTransaction{
		UserID:          userID,
		TransactionType: "redemption",
		Amount:          -product.PointsRequired,
		BalanceAfter:    newBalance,
		Reason:          fmt.Sprintf("兑换商品: %s / Redeem product: %s", product.Name, product.Name),
		OperatorID:      nil, // User operation
		RelatedOrderID:  &order.ID,
	}

	err = tx.Create(transaction).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	// Load relationships for response
	order.User = user
	order.Product = product

	return order, nil
}

// generateOrderNumber generates a unique order number
// Format: RD + timestamp + userID + productID
func (s *RedemptionService) generateOrderNumber(userID, productID uint) string {
	timestamp := time.Now().Format("20060102150405")
	return fmt.Sprintf("RD%s%d%d", timestamp, userID, productID)
}

// GetRedemptionHistory gets a user's redemption history
func (s *RedemptionService) GetRedemptionHistory(userID uint) ([]models.RedemptionOrder, error) {
	return s.orderRepo.GetByUserID(userID)
}

// GetOrderByID gets an order by ID
func (s *RedemptionService) GetOrderByID(orderID uint) (*models.RedemptionOrder, error) {
	return s.orderRepo.GetByID(orderID)
}

// GetOrderByNumber gets an order by order number
func (s *RedemptionService) GetOrderByNumber(orderNumber string) (*models.RedemptionOrder, error) {
	return s.orderRepo.GetByOrderNumber(orderNumber)
}

// UpdateOrderStatus updates an order's status
func (s *RedemptionService) UpdateOrderStatus(orderID uint, status string) error {
	if status != "preparing" && status != "delivered" {
		return errors.New("invalid status: must be 'preparing' or 'delivered'")
	}

	return s.orderRepo.UpdateStatus(orderID, status)
}

// BatchUpdateOrderStatus updates multiple orders' status
func (s *RedemptionService) BatchUpdateOrderStatus(orderNumbers []string, status string) error {
	if status != "preparing" && status != "delivered" {
		return errors.New("invalid status: must be 'preparing' or 'delivered'")
	}

	if len(orderNumbers) == 0 {
		return errors.New("order numbers list is empty")
	}

	// Validate all order numbers exist
	for _, orderNumber := range orderNumbers {
		exists, err := s.orderRepo.ExistsByOrderNumber(orderNumber)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("order %s not found", orderNumber)
		}
	}

	return s.orderRepo.BatchUpdateStatus(orderNumbers, status)
}

// ParseBatchOrderNumbers parses a list of order numbers from various formats
// Supports: comma-separated, newline-separated, or space-separated
func (s *RedemptionService) ParseBatchOrderNumbers(input string) []string {
	input = strings.TrimSpace(input)
	
	// Try different separators
	var orderNumbers []string
	
	// Check for newlines first
	if strings.Contains(input, "\n") {
		lines := strings.Split(input, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				orderNumbers = append(orderNumbers, line)
			}
		}
	} else if strings.Contains(input, ",") {
		// Comma-separated
		parts := strings.Split(input, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				orderNumbers = append(orderNumbers, part)
			}
		}
	} else {
		// Space-separated or single order
		parts := strings.Fields(input)
		orderNumbers = append(orderNumbers, parts...)
	}

	return orderNumbers
}

// GetAllOrders gets all orders (for admin)
func (s *RedemptionService) GetAllOrders() ([]models.RedemptionOrder, error) {
	return s.orderRepo.GetAll()
}

// ListOrders lists orders with optional filters
func (s *RedemptionService) ListOrders(status *string, userID *uint) ([]models.RedemptionOrder, error) {
	return s.orderRepo.List(status, userID)
}

// GetRedemptionStats gets redemption statistics for reports
func (s *RedemptionService) GetRedemptionStats() ([]repository.OrderStats, error) {
	return s.orderRepo.GetOrderStats()
}

// CountOrdersByStatus counts orders by status
func (s *RedemptionService) CountOrdersByStatus(status string) (int64, error) {
	if status != "preparing" && status != "delivered" {
		return 0, errors.New("invalid status")
	}
	return s.orderRepo.CountByStatus(status)
}

// ValidateRedemption validates if a redemption can proceed
func (s *RedemptionService) ValidateRedemption(userID uint, productID uint) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if !user.IsActive {
		return errors.New("user account is inactive")
	}

	// Get product
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return errors.New("product not found")
	}

	if product.Status != "active" {
		return errors.New("product is not available")
	}

	if product.StockQuantity <= 0 {
		return errors.New("product is out of stock")
	}

	if user.PointsBalance < product.PointsRequired {
		return errors.New("insufficient points")
	}

	return nil
}
