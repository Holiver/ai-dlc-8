package repository

import (
	"awsome-shop/internal/models"
	"errors"

	"gorm.io/gorm"
)

// PointsTransactionRepository handles points transaction data access operations
type PointsTransactionRepository struct {
	db *gorm.DB
}

// NewPointsTransactionRepository creates a new PointsTransactionRepository instance
func NewPointsTransactionRepository(db *gorm.DB) *PointsTransactionRepository {
	return &PointsTransactionRepository{db: db}
}

// Create creates a new points transaction
func (r *PointsTransactionRepository) Create(transaction *models.PointsTransaction) error {
	return r.db.Create(transaction).Error
}

// GetByID retrieves a points transaction by ID
func (r *PointsTransactionRepository) GetByID(id uint) (*models.PointsTransaction, error) {
	var transaction models.PointsTransaction
	err := r.db.Preload("User").Preload("Operator").Preload("RelatedOrder").
		First(&transaction, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

// GetByUserID retrieves all transactions for a specific user (ordered by created_at desc)
func (r *PointsTransactionRepository) GetByUserID(userID uint) ([]models.PointsTransaction, error) {
	var transactions []models.PointsTransaction
	err := r.db.Preload("Operator").Preload("RelatedOrder").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

// GetByUserIDWithPagination retrieves transactions for a user with pagination
func (r *PointsTransactionRepository) GetByUserIDWithPagination(userID uint, page, pageSize int) ([]models.PointsTransaction, int64, error) {
	var transactions []models.PointsTransaction
	var total int64

	// Count total records
	err := r.db.Model(&models.PointsTransaction{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated records
	offset := (page - 1) * pageSize
	err = r.db.Preload("Operator").Preload("RelatedOrder").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&transactions).Error

	return transactions, total, err
}

// List retrieves all transactions with optional filters
func (r *PointsTransactionRepository) List(userID *uint, transactionType *string) ([]models.PointsTransaction, error) {
	var transactions []models.PointsTransaction
	query := r.db.Preload("User").Preload("Operator").Preload("RelatedOrder")

	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}

	if transactionType != nil {
		query = query.Where("transaction_type = ?", *transactionType)
	}

	err := query.Order("created_at DESC").Find(&transactions).Error
	return transactions, err
}

// GetGrantTransactions retrieves all grant transactions for reports
type GrantTransactionStats struct {
	UserName     string
	UserEmail    string
	Amount       int
	Reason       string
	OperatorName string
	CreatedAt    string
}

func (r *PointsTransactionRepository) GetGrantTransactions() ([]GrantTransactionStats, error) {
	var stats []GrantTransactionStats
	err := r.db.Table("points_transactions").
		Select("users.full_name as user_name, users.email as user_email, points_transactions.amount, points_transactions.reason, operators.full_name as operator_name, points_transactions.created_at").
		Joins("LEFT JOIN users ON users.id = points_transactions.user_id").
		Joins("LEFT JOIN users as operators ON operators.id = points_transactions.operator_id").
		Where("points_transactions.transaction_type = ?", "grant").
		Order("points_transactions.created_at DESC").
		Scan(&stats).Error
	return stats, err
}

// GetPointsBalances retrieves current points balances for all users
type PointsBalanceStats struct {
	UserName      string
	UserEmail     string
	PointsBalance int
}

func (r *PointsTransactionRepository) GetPointsBalances() ([]PointsBalanceStats, error) {
	var stats []PointsBalanceStats
	err := r.db.Table("users").
		Select("users.full_name as user_name, users.email as user_email, users.points_balance").
		Where("users.is_active = ?", true).
		Order("users.points_balance DESC").
		Scan(&stats).Error
	return stats, err
}

// CountByType counts transactions by type
func (r *PointsTransactionRepository) CountByType(transactionType string) (int64, error) {
	var count int64
	err := r.db.Model(&models.PointsTransaction{}).
		Where("transaction_type = ?", transactionType).
		Count(&count).Error
	return count, err
}

// GetTotalPointsGranted calculates total points granted
func (r *PointsTransactionRepository) GetTotalPointsGranted() (int, error) {
	var total int
	err := r.db.Model(&models.PointsTransaction{}).
		Where("transaction_type = ?", "grant").
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

// GetTotalPointsRedeemed calculates total points redeemed
func (r *PointsTransactionRepository) GetTotalPointsRedeemed() (int, error) {
	var total int
	err := r.db.Model(&models.PointsTransaction{}).
		Where("transaction_type = ?", "redemption").
		Select("COALESCE(SUM(ABS(amount)), 0)").
		Scan(&total).Error
	return total, err
}

// BatchCreate creates multiple transactions in a single operation
func (r *PointsTransactionRepository) BatchCreate(transactions []models.PointsTransaction) error {
	return r.db.Create(&transactions).Error
}

// GetByUserIDAndType retrieves transactions for a user filtered by type
func (r *PointsTransactionRepository) GetByUserIDAndType(userID uint, transactionType string) ([]models.PointsTransaction, error) {
	var transactions []models.PointsTransaction
	err := r.db.Preload("Operator").Preload("RelatedOrder").
		Where("user_id = ? AND transaction_type = ?", userID, transactionType).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}
