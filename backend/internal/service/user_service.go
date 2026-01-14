package service

import (
	"awsome-shop/internal/models"
	"awsome-shop/internal/repository"
	"errors"
	"fmt"
	"regexp"

	"gorm.io/gorm"
)

// UserService handles user management operations
type UserService struct {
	userRepo            *repository.UserRepository
	pointsTransactionRepo *repository.PointsTransactionRepository
	authService         *AuthService
	db                  *gorm.DB
}

// NewUserService creates a new UserService instance
func NewUserService(
	userRepo *repository.UserRepository,
	pointsTransactionRepo *repository.PointsTransactionRepository,
	authService *AuthService,
	db *gorm.DB,
) *UserService {
	return &UserService{
		userRepo:            userRepo,
		pointsTransactionRepo: pointsTransactionRepo,
		authService:         authService,
		db:                  db,
	}
}

// CreateEmployeeRequest represents a request to create an employee
type CreateEmployeeRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
}

// CreateEmployee creates a new employee account
// Initial password is set to the last 6 digits of the phone number
func (s *UserService) CreateEmployee(req *CreateEmployeeRequest) (*models.User, error) {
	// Validate phone number format
	if len(req.Phone) < 6 {
		return nil, errors.New("phone number must be at least 6 digits")
	}

	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	// Generate initial password (last 6 digits of phone)
	initialPassword := req.Phone[len(req.Phone)-6:]

	// Hash password
	hashedPassword, err := s.authService.HashPassword(initialPassword)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		FullName:          req.FullName,
		Email:             req.Email,
		Phone:             req.Phone,
		PasswordHash:      hashedPassword,
		Role:              "employee",
		PointsBalance:     0,
		IsFirstLogin:      true,
		IsActive:          true,
		PreferredLanguage: "zh",
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""

	return user, nil
}

// UpdatePhoneRequest represents a request to update phone number
type UpdatePhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// UpdatePhone updates a user's phone number
func (s *UserService) UpdatePhone(userID uint, newPhone string) error {
	// Validate phone number
	if len(newPhone) < 6 {
		return errors.New("phone number must be at least 6 digits")
	}

	// Validate phone format (basic validation)
	matched, _ := regexp.MatchString(`^\d+$`, newPhone)
	if !matched {
		return errors.New("phone number must contain only digits")
	}

	// Update phone
	err := s.userRepo.UpdatePhone(userID, newPhone)
	if err != nil {
		return err
	}

	return nil
}

// SetEmployeeDeparture sets an employee's status to inactive (离职)
// This also invalidates their points
func (s *UserService) SetEmployeeDeparture(userID uint, operatorID uint) error {
	// Get user
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	// Check if user is already inactive
	if !user.IsActive {
		return errors.New("user is already inactive")
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

	// Set user as inactive
	err = tx.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", false).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// If user has points, create a deduction transaction to invalidate them
	if user.PointsBalance > 0 {
		transaction := &models.PointsTransaction{
			UserID:          userID,
			TransactionType: "deduct",
			Amount:          -user.PointsBalance,
			BalanceAfter:    0,
			Reason:          fmt.Sprintf("员工离职，积分失效 / Employee departure, points invalidated"),
			OperatorID:      &operatorID,
			RelatedOrderID:  nil,
		}

		err = tx.Create(transaction).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// Set points balance to 0
		err = tx.Model(&models.User{}).
			Where("id = ?", userID).
			Update("points_balance", 0).Error
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

// GetUserProfile gets a user's profile information
func (s *UserService) GetUserProfile(userID uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""

	return user, nil
}

// GetUserByEmail gets a user by email
func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	// Don't return password hash
	user.PasswordHash = ""

	return user, nil
}

// ListUsers lists all users with optional filters
func (s *UserService) ListUsers(isActive *bool) ([]models.User, error) {
	users, err := s.userRepo.List(isActive)
	if err != nil {
		return nil, err
	}

	// Don't return password hashes
	for i := range users {
		users[i].PasswordHash = ""
	}

	return users, nil
}

// IsAdmin checks if a user is an administrator
func (s *UserService) IsAdmin(userID uint) (bool, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	return user.Role == "admin", nil
}

// ValidateUserActive checks if a user is active
func (s *UserService) ValidateUserActive(userID uint) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if !user.IsActive {
		return errors.New("user account is inactive")
	}

	return nil
}
