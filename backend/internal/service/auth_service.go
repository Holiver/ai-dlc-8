package service

import (
	"awsome-shop/internal/models"
	"awsome-shop/internal/repository"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService handles authentication operations
type AuthService struct {
	userRepo            *repository.UserRepository
	pointsTransactionRepo *repository.PointsTransactionRepository
	db                  *gorm.DB
	jwtSecret           string
	jwtExpirationHours  int
}

// NewAuthService creates a new AuthService instance
func NewAuthService(
	userRepo *repository.UserRepository,
	pointsTransactionRepo *repository.PointsTransactionRepository,
	db *gorm.DB,
	jwtSecret string,
	jwtExpirationHours int,
) *AuthService {
	return &AuthService{
		userRepo:            userRepo,
		pointsTransactionRepo: pointsTransactionRepo,
		db:                  db,
		jwtSecret:           jwtSecret,
		jwtExpirationHours:  jwtExpirationHours,
	}
}

// Claims represents JWT claims
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// HashPassword hashes a password using bcrypt
func (s *AuthService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// VerifyPassword verifies a password against a hash
func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// GenerateToken generates a JWT token for a user
func (s *AuthService) GenerateToken(user *models.User) (string, error) {
	expirationTime := time.Now().Add(time.Duration(s.jwtExpirationHours) * time.Hour)

	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *AuthService) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

// Login authenticates a user and returns a JWT token
// Handles first-time login by granting 1000 initial points
func (s *AuthService) Login(email, password string) (*LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	err = s.VerifyPassword(user.PasswordHash, password)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Handle first-time login
	if user.IsFirstLogin {
		err = s.handleFirstLogin(user)
		if err != nil {
			return nil, err
		}
	}

	// Generate JWT token
	token, err := s.GenerateToken(user)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	// Don't return password hash
	user.PasswordHash = ""

	return &LoginResponse{
		Token: token,
		User:  user,
	}, nil
}

// handleFirstLogin handles first-time login by granting initial points
func (s *AuthService) handleFirstLogin(user *models.User) error {
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

	// Grant 1000 initial points
	initialPoints := 1000
	newBalance := user.PointsBalance + initialPoints

	// Update user points balance
	err := tx.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"points_balance":  newBalance,
			"is_first_login": false,
		}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// Create points transaction record
	transaction := &models.PointsTransaction{
		UserID:          user.ID,
		TransactionType: "grant",
		Amount:          initialPoints,
		BalanceAfter:    newBalance,
		Reason:          "首次登录奖励 / First login bonus",
		OperatorID:      nil, // System operation
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

	// Update user object
	user.PointsBalance = newBalance
	user.IsFirstLogin = false

	return nil
}

// GetUserFromToken extracts user information from a JWT token
func (s *AuthService) GetUserFromToken(tokenString string) (*models.User, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Don't return password hash
	user.PasswordHash = ""

	return user, nil
}
