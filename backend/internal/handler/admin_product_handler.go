package handler

import (
	"awsome-shop/internal/middleware"
	"awsome-shop/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminProductHandler handles admin product management requests
type AdminProductHandler struct {
	productService *service.ProductService
}

// NewAdminProductHandler creates a new AdminProductHandler instance
func NewAdminProductHandler(productService *service.ProductService) *AdminProductHandler {
	return &AdminProductHandler{
		productService: productService,
	}
}

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	Name           string `json:"name" binding:"required"`
	ImageURL       string `json:"image_url"`
	PointsRequired int    `json:"points_required" binding:"required,min=1"`
	StockQuantity  int    `json:"stock_quantity" binding:"min=0"`
}

// CreateProduct creates a new product
// POST /api/v1/admin/products
func (h *AdminProductHandler) CreateProduct(c *gin.Context) {
	operatorID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Operator ID not found in context",
		})
		return
	}

	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	createReq := &service.CreateProductRequest{
		Name:           req.Name,
		ImageURL:       req.ImageURL,
		PointsRequired: req.PointsRequired,
		StockQuantity:  req.StockQuantity,
	}

	product, err := h.productService.CreateProduct(createReq, operatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"product": product,
	})
}

// UpdateProductRequest represents a request to update a product
type UpdateProductRequest struct {
	Name           *string `json:"name"`
	ImageURL       *string `json:"image_url"`
	PointsRequired *int    `json:"points_required"`
	StockQuantity  *int    `json:"stock_quantity"`
}

// UpdateProduct updates a product
// PUT /api/v1/admin/products/:id
func (h *AdminProductHandler) UpdateProduct(c *gin.Context) {
	operatorID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Operator ID not found in context",
		})
		return
	}

	// Get product ID from URL
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	updateReq := &service.UpdateProductRequest{
		Name:           req.Name,
		ImageURL:       req.ImageURL,
		PointsRequired: req.PointsRequired,
		StockQuantity:  req.StockQuantity,
	}

	product, err := h.productService.UpdateProduct(uint(productID), updateReq, operatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// SetProductStatusRequest represents a request to set product status
type SetProductStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}

// SetProductStatus sets a product's status (active/inactive)
// PUT /api/v1/admin/products/:id/status
func (h *AdminProductHandler) SetProductStatus(c *gin.Context) {
	// Get product ID from URL
	productIDStr := c.Param("id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	var req SetProductStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format. Status must be 'active' or 'inactive'",
		})
		return
	}

	err = h.productService.SetProductStatus(uint(productID), req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product status updated successfully",
	})
}

// BatchImportProductsRequest represents a request to batch import products
type BatchImportProductsRequest struct {
	Markdown string `json:"markdown" binding:"required"`
}

// BatchImportProducts imports multiple products from markdown table
// POST /api/v1/admin/products/batch
func (h *AdminProductHandler) BatchImportProducts(c *gin.Context) {
	operatorID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Operator ID not found in context",
		})
		return
	}

	var req BatchImportProductsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	products, err := h.productService.BatchImportProducts(req.Markdown, operatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"products": products,
		"count":    len(products),
		"message":  "Products imported successfully",
	})
}

// ListAllProducts lists all products (including inactive)
// GET /api/v1/admin/products
func (h *AdminProductHandler) ListAllProducts(c *gin.Context) {
	// Get optional status filter
	var status *string
	if statusStr := c.Query("status"); statusStr != "" {
		status = &statusStr
	}

	products, err := h.productService.ListProducts(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

// RegisterRoutes registers admin product routes
func (h *AdminProductHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware, adminMiddleware gin.HandlerFunc) {
	admin := router.Group("/admin/products")
	admin.Use(authMiddleware, adminMiddleware)
	{
		admin.POST("", h.CreateProduct)
		admin.PUT("/:id", h.UpdateProduct)
		admin.PUT("/:id/status", h.SetProductStatus)
		admin.POST("/batch", h.BatchImportProducts)
		admin.GET("", h.ListAllProducts)
	}
}
