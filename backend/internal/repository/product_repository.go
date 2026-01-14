package repository

import (
	"awsome-shop/internal/models"
	"errors"

	"gorm.io/gorm"
)

// ProductRepository handles product data access operations
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository instance
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product
func (r *ProductRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

// GetByID retrieves a product by ID
func (r *ProductRepository) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// Update updates a product
func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

// Delete deletes a product
func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

// GetActiveProducts retrieves all active (上架) products
func (r *ProductRepository) GetActiveProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("status = ?", "active").Find(&products).Error
	return products, err
}

// List retrieves all products with optional status filter
func (r *ProductRepository) List(status *string) ([]models.Product, error) {
	var products []models.Product
	query := r.db

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	err := query.Find(&products).Error
	return products, err
}

// UpdateStock updates a product's stock quantity
func (r *ProductRepository) UpdateStock(productID uint, newStock int) error {
	return r.db.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("stock_quantity", newStock).Error
}

// DecrementStock decrements a product's stock by the specified amount
// Returns error if stock is insufficient
func (r *ProductRepository) DecrementStock(productID uint, amount int) error {
	result := r.db.Model(&models.Product{}).
		Where("id = ? AND stock_quantity >= ?", productID, amount).
		Update("stock_quantity", gorm.Expr("stock_quantity - ?", amount))

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("insufficient stock or product not found")
	}

	return nil
}

// UpdateStatus updates a product's status (active/inactive)
func (r *ProductRepository) UpdateStatus(productID uint, status string) error {
	if status != "active" && status != "inactive" {
		return errors.New("invalid status: must be 'active' or 'inactive'")
	}

	return r.db.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("status", status).Error
}

// UpdatePoints updates a product's required points
func (r *ProductRepository) UpdatePoints(productID uint, points int) error {
	return r.db.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("points_required", points).Error
}

// GetByIDWithLock retrieves a product by ID with row lock (for transaction)
func (r *ProductRepository) GetByIDWithLock(tx *gorm.DB, id uint) (*models.Product, error) {
	var product models.Product
	err := tx.Clauses(gorm.Locking{Strength: "UPDATE"}).First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// BatchCreate creates multiple products in a single transaction
func (r *ProductRepository) BatchCreate(products []models.Product) error {
	return r.db.Create(&products).Error
}
