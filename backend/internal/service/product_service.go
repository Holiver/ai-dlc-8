package service

import (
	"awsome-shop/internal/models"
	"awsome-shop/internal/repository"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// ProductService handles product management operations
type ProductService struct {
	productRepo        *repository.ProductRepository
	priceHistoryRepo   *repository.ProductRepository // For price history operations
	db                 *gorm.DB
}

// NewProductService creates a new ProductService instance
func NewProductService(
	productRepo *repository.ProductRepository,
	db *gorm.DB,
) *ProductService {
	return &ProductService{
		productRepo:      productRepo,
		priceHistoryRepo: productRepo,
		db:               db,
	}
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	Name           string `json:"name" binding:"required"`
	ImageURL       string `json:"image_url"`
	PointsRequired int    `json:"points_required" binding:"required,min=1"`
	StockQuantity  int    `json:"stock_quantity" binding:"min=0"`
}

// CreateProduct creates a new product and records initial price history
func (s *ProductService) CreateProduct(req *CreateProductRequest, operatorID uint) (*models.Product, error) {
	// Validate points
	if req.PointsRequired <= 0 {
		return nil, errors.New("points required must be greater than 0")
	}

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

	// Create product
	product := &models.Product{
		Name:           req.Name,
		ImageURL:       req.ImageURL,
		PointsRequired: req.PointsRequired,
		StockQuantity:  req.StockQuantity,
		Status:         "active",
	}

	err := tx.Create(product).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create price history record
	priceHistory := &models.ProductPriceHistory{
		ProductID:  product.ID,
		OldPoints:  nil, // No old price for new product
		NewPoints:  req.PointsRequired,
		OperatorID: &operatorID,
	}

	err = tx.Create(priceHistory).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Name           *string `json:"name"`
	ImageURL       *string `json:"image_url"`
	PointsRequired *int    `json:"points_required"`
	StockQuantity  *int    `json:"stock_quantity"`
}

// UpdateProduct updates a product and records price change if applicable
func (s *ProductService) UpdateProduct(productID uint, req *UpdateProductRequest, operatorID uint) (*models.Product, error) {
	// Get existing product
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return nil, err
	}

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

	// Track if points changed
	pointsChanged := false
	oldPoints := product.PointsRequired

	// Update fields
	if req.Name != nil {
		product.Name = *req.Name
	}
	if req.ImageURL != nil {
		product.ImageURL = *req.ImageURL
	}
	if req.PointsRequired != nil {
		if *req.PointsRequired <= 0 {
			tx.Rollback()
			return nil, errors.New("points required must be greater than 0")
		}
		if *req.PointsRequired != product.PointsRequired {
			pointsChanged = true
			product.PointsRequired = *req.PointsRequired
		}
	}
	if req.StockQuantity != nil {
		if *req.StockQuantity < 0 {
			tx.Rollback()
			return nil, errors.New("stock quantity cannot be negative")
		}
		product.StockQuantity = *req.StockQuantity
	}

	// Save product
	err = tx.Save(product).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// If points changed, create price history record
	if pointsChanged {
		priceHistory := &models.ProductPriceHistory{
			ProductID:  product.ID,
			OldPoints:  &oldPoints,
			NewPoints:  product.PointsRequired,
			OperatorID: &operatorID,
		}

		err = tx.Create(priceHistory).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

// SetProductStatus sets a product's status (active/inactive)
func (s *ProductService) SetProductStatus(productID uint, status string) error {
	if status != "active" && status != "inactive" {
		return errors.New("invalid status: must be 'active' or 'inactive'")
	}

	return s.productRepo.UpdateStatus(productID, status)
}

// GetActiveProducts retrieves all active products
func (s *ProductService) GetActiveProducts() ([]models.Product, error) {
	return s.productRepo.GetActiveProducts()
}

// GetProductByID retrieves a product by ID
func (s *ProductService) GetProductByID(productID uint) (*models.Product, error) {
	return s.productRepo.GetByID(productID)
}

// ListProducts lists all products with optional status filter
func (s *ProductService) ListProducts(status *string) ([]models.Product, error) {
	return s.productRepo.List(status)
}

// BatchImportProduct represents a product for batch import
type BatchImportProduct struct {
	Name          string
	ImageURL      string
	StockQuantity int
	PointsRequired int
}

// ParseMarkdownTable parses a markdown table for batch product import
// Expected format:
// | 商品名称（中英文） | 商品主图 | 商品数量 | 所需积分 |
// |-------------------|---------|---------|---------|
// | Product Name      | url     | 10      | 100     |
func (s *ProductService) ParseMarkdownTable(markdown string) ([]BatchImportProduct, error) {
	lines := strings.Split(strings.TrimSpace(markdown), "\n")
	
	if len(lines) < 3 {
		return nil, errors.New("invalid markdown table: must have header, separator, and at least one data row")
	}

	// Skip header (line 0) and separator (line 1)
	var products []BatchImportProduct
	
	for i := 2; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}

		// Parse table row
		// Remove leading and trailing pipes
		line = strings.Trim(line, "|")
		
		// Split by pipe
		fields := strings.Split(line, "|")
		if len(fields) < 4 {
			return nil, fmt.Errorf("invalid row %d: expected 4 columns", i-1)
		}

		// Trim whitespace from each field
		for j := range fields {
			fields[j] = strings.TrimSpace(fields[j])
		}

		// Parse stock quantity
		stockQuantity, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, fmt.Errorf("invalid stock quantity in row %d: %s", i-1, fields[2])
		}

		// Parse points required (optional, default to 0 if not provided or in wrong format)
		pointsRequired := 0
		if len(fields) > 3 && fields[3] != "" {
			pointsRequired, err = strconv.Atoi(fields[3])
			if err != nil {
				return nil, fmt.Errorf("invalid points required in row %d: %s", i-1, fields[3])
			}
		}

		product := BatchImportProduct{
			Name:          fields[0],
			ImageURL:      fields[1],
			StockQuantity: stockQuantity,
			PointsRequired: pointsRequired,
		}

		products = append(products, product)
	}

	if len(products) == 0 {
		return nil, errors.New("no valid products found in markdown table")
	}

	return products, nil
}

// BatchImportProducts imports multiple products from markdown table
func (s *ProductService) BatchImportProducts(markdown string, operatorID uint) ([]models.Product, error) {
	// Parse markdown table
	importProducts, err := s.ParseMarkdownTable(markdown)
	if err != nil {
		return nil, err
	}

	// Validate all products first
	for i, p := range importProducts {
		if p.Name == "" {
			return nil, fmt.Errorf("product %d: name is required", i+1)
		}
		if p.StockQuantity < 0 {
			return nil, fmt.Errorf("product %d: stock quantity cannot be negative", i+1)
		}
		if p.PointsRequired <= 0 {
			return nil, fmt.Errorf("product %d: points required must be greater than 0", i+1)
		}
	}

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

	var createdProducts []models.Product

	// Create all products
	for _, p := range importProducts {
		product := models.Product{
			Name:           p.Name,
			ImageURL:       p.ImageURL,
			PointsRequired: p.PointsRequired,
			StockQuantity:  p.StockQuantity,
			Status:         "active",
		}

		err = tx.Create(&product).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		// Create price history
		priceHistory := &models.ProductPriceHistory{
			ProductID:  product.ID,
			OldPoints:  nil,
			NewPoints:  product.PointsRequired,
			OperatorID: &operatorID,
		}

		err = tx.Create(priceHistory).Error
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		createdProducts = append(createdProducts, product)
	}

	// Commit transaction
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return createdProducts, nil
}

// ValidateProductAvailable checks if a product is available for redemption
func (s *ProductService) ValidateProductAvailable(productID uint) error {
	product, err := s.productRepo.GetByID(productID)
	if err != nil {
		return err
	}

	if product.Status != "active" {
		return errors.New("product is not active")
	}

	if product.StockQuantity <= 0 {
		return errors.New("product is out of stock")
	}

	return nil
}

// IsValidMarkdownTable validates if a string is a valid markdown table format
func (s *ProductService) IsValidMarkdownTable(markdown string) bool {
	lines := strings.Split(strings.TrimSpace(markdown), "\n")
	
	if len(lines) < 3 {
		return false
	}

	// Check if second line is a separator (contains dashes and pipes)
	separatorPattern := regexp.MustCompile(`^[\s\|:\-]+$`)
	return separatorPattern.MatchString(lines[1])
}
