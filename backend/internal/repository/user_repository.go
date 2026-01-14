package repository

import (
	"awsome-shop/internal/models"
	"errors"

	"gorm.io/gorm"
)

// UserRepository handles user data access operations
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete deletes a user (soft delete if configured)
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// UpdatePointsBalance updates a user's points balance
func (r *UserRepository) UpdatePointsBalance(userID uint, newBalance int) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("points_balance", newBalance).Error
}

// UpdateStatus updates a user's active status (for handling employee departure)
func (r *UserRepository) UpdateStatus(userID uint, isActive bool) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_active", isActive).Error
}

// UpdatePhone updates a user's phone number
func (r *UserRepository) UpdatePhone(userID uint, phone string) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("phone", phone).Error
}

// UpdateFirstLoginFlag updates the is_first_login flag
func (r *UserRepository) UpdateFirstLoginFlag(userID uint, isFirstLogin bool) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("is_first_login", isFirstLogin).Error
}

// List retrieves all users with optional filters
func (r *UserRepository) List(isActive *bool) ([]models.User, error) {
	var users []models.User
	query := r.db

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	err := query.Find(&users).Error
	return users, err
}

// GetActiveUsers retrieves all active users
func (r *UserRepository) GetActiveUsers() ([]models.User, error) {
	var users []models.User
	err := r.db.Where("is_active = ?", true).Find(&users).Error
	return users, err
}

// ExistsByEmail checks if a user with the given email exists
func (r *UserRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}
