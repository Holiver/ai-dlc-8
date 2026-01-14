package service

import (
	"awsome-shop/internal/models"
	"awsome-shop/internal/repository"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// PointsService handles points management operations
type PointsService struct {
	userRepo              *repository.UserRepository
	pointsTransactionRepo *repository.PointsTransactionRepository
	db                    *gorm.DB
}

// NewPointsService creates a new PointsService instance
func NewPointsService(
	userRepo *repository.UserRepository,
	pointsTransactionRepo *repository.PointsTransactionRepository,
	db *gorm.DB,
) *PointsService {
	return &PointsService{
		userRepo:              userRepo,
		pointsTransactionRepo: pointsTransactionRepo,
		db:                    db,
	}
}

// GrantPointsRequest represents a request to grant points
type GrantPointsRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Amount int    `json:"amount" binding:"required,min=1"`
	Reason string `json:"reason" binding:"required"`
}

// GrantPoints grants points to a user
func (s *PointsService) GrantPoints(userID uint, amount int, reason string, operatorID uint) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if reason == "" {
		return errors.New("reason is required")
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return errors.New("cannot grant points to inactive user")
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate new balance
	newBalance := user.PointsBalance + amount

	// Update user points balance
	err = tx.Model(&models.User{}).
		Where("id = ?", userID).
		Update("points_balance", newBalance).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create transaction record
	transaction := &models.PointsTransaction{
		UserID:          userID,
		TransactionType: "grant",
		Amount:          amount,
		BalanceAfter:    newBalance,
		Reason:          reason,
		OperatorID:      &operatorID,
		RelatedOrderID:  nil,
	}

	err = tx.Create(transaction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

// DeductPointsRequest represents a request to deduct points
type DeductPointsRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Amount int    `json:"amount" binding:"required,min=1"`
	Reason string `json:"reason" binding:"required"`
}

// DeductPoints deducts points from a user
func (s *PointsService) DeductPoints(userID uint, amount int, reason string, operatorID uint) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	if reason == "" {
		return errors.New("reason is required")
	}

	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Check if user has sufficient balance
	if user.PointsBalance < amount {
		return errors.New("insufficient points balance")
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate new balance
	newBalance := user.PointsBalance - amount

	// Update user points balance
	err = tx.Model(&models.User{}).
		Where("id = ?", userID).
		Update("points_balance", newBalance).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create transaction record (negative amount for deduction)
	transaction := &models.PointsTransaction{
		UserID:          userID,
		TransactionType: "deduct",
		Amount:          -amount,
		BalanceAfter:    newBalance,
		Reason:          reason,
		OperatorID:      &operatorID,
		RelatedOrderID:  nil,
	}

	err = tx.Create(transaction).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

// GetPointsBalance gets a user's current points balance
func (s *PointsService) GetPointsBalance(userID uint) (int, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return 0, err
	}

	return user.PointsBalance, nil
}

// GetPointsHistory gets a user's points transaction history with pagination
func (s *PointsService) GetPointsHistory(userID uint, page, pageSize int) ([]models.PointsTransaction, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	return s.pointsTransactionRepo.GetByUserIDWithPagination(userID, page, pageSize)
}

// BatchGrantEntry represents a single entry in batch grant operation
type BatchGrantEntry struct {
	Email  string
	Name   string
	Amount int
	Reason string
}

// ParseBatchGrantMarkdown parses markdown table for batch points grant
// Expected format:
// | 员工邮箱 | 姓名 | 积分 | 备注 |
// |---------|------|------|------|
// | email   | name | 100  | note |
func (s *PointsService) ParseBatchGrantMarkdown(markdown string) ([]BatchGrantEntry, error) {
	lines := strings.Split(strings.TrimSpace(markdown), "\n")
	
	if len(lines) < 3 {
		return nil, errors.New("invalid markdown table: must have header, separator, and at least one data row")
	}

	var entries []BatchGrantEntry
	
	// Skip header (line 0) and separator (line 1)
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Parse table row
		line = strings.Trim(line, "|")
		fields := strings.Split(line, "|")
		
		if len(fields) < 4 {
			return nil, fmt.Errorf("invalid row %d: expected 4 columns (email, name, amount, reason)", i-1)
		}

		// Trim whitespace
		for j := range fields {
			fields[j] = strings.TrimSpace(fields[j])
		}

		// Parse amount
		amount, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("invalid amount in row %d: %s", i-1, fields[2])
		}

		entry := BatchGrantEntry{
			Email:  fields[0],
			Name:   fields[1],
			Amount: amount,
			Reason: fields[3],
		}

		entries = append(entries, entry)
	}

	if len(entries) == 0 {
		return nil, errors.New("no valid entries found in markdown table")
	}

	return entries, nil
}

// BatchGrantPoints grants points to multiple users
func (s *PointsService) BatchGrantPoints(markdown string, operatorID uint) error {
	// Parse markdown table
	entries, err := s.ParseBatchGrantMarkdown(markdown)
	if err != nil {
		return err
	}

	// Validate all entries first
	userMap := make(map[string]*models.User)
	
	for i, entry := range entries {
		if entry.Email == "" {
			return fmt.Errorf("entry %d: email is required", i+1)
		}
		if entry.Amount <= 0 {
			return fmt.Errorf("entry %d: amount must be greater than 0", i+1)
		}
		if entry.Reason == "" {
			return fmt.Errorf("entry %d: reason is required", i+1)
		}

		// Get user by email
		user, err := s.userRepo.GetByEmail(entry.Email)
		if err != nil {
			return fmt.Errorf("entry %d: user with email %s not found", i+1, entry.Email)
		}

		if !user.IsActive {
			return fmt.Errorf("entry %d: user %s is inactive", i+1, entry.Email)
		}

		// Verify name matches (optional check)
		if entry.Name != "" && entry.Name != user.FullName {
			return fmt.Errorf("entry %d: name mismatch for %s (expected: %s, got: %s)", 
				i+1, entry.Email, user.FullName, entry.Name)
		}

		userMap[entry.Email] = user
	}

	// Start transaction
	tx := s.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Process all entries
	for _, entry := range entries {
		user := userMap[entry.Email]
		newBalance := user.PointsBalance + entry.Amount

		// Update user points balance
		err = tx.Model(&models.User{}).
			Where("id = ?", user.ID).
			Update("points_balance", newBalance).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create transaction record
		transaction := &models.PointsTransaction{
			UserID:          user.ID,
			TransactionType: "grant",
			Amount:          entry.Amount,
			BalanceAfter:    newBalance,
			Reason:          entry.Reason,
			OperatorID:      &operatorID,
			RelatedOrderID:  nil,
		}

		err = tx.Create(transaction).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return err
	}

	return nil
}

// GetGrantTransactionsReport gets points grant transactions for reporting
func (s *PointsService) GetGrantTransactionsReport() ([]repository.GrantTransactionStats, error) {
	return s.pointsTransactionRepo.GetGrantTransactions()
}

// GetPointsBalancesReport gets current points balances for all users
func (s *PointsService) GetPointsBalancesReport() ([]repository.PointsBalanceStats, error) {
	return s.pointsTransactionRepo.GetPointsBalances()
}

// GetTotalPointsGranted gets total points granted in the system
func (s *PointsService) GetTotalPointsGranted() (int, error) {
	return s.pointsTransactionRepo.GetTotalPointsGranted()
}

// GetTotalPointsRedeemed gets total points redeemed in the system
func (s *PointsService) GetTotalPointsRedeemed() (int, error) {
	return s.pointsTransactionRepo.GetTotalPointsRedeemed()
}
