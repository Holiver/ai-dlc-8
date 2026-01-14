package handler

import (
	"awsome-shop/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ProductHandler handles product related requests
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new ProductHandler instance
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// GetProducts gets all active products
// GET /api/v1/products
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productService.GetActiveProducts()
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

// GetProductByID gets a product by ID
// GET /api/v1/products/:id
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, err := h.productService.GetProductByID(uri.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

// RegisterRoutes registers product routes
func (h *ProductHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	products := router.Group("/products")
	products.Use(authMiddleware)
	{
		products.GET("", h.GetProducts)
		products.GET("/:id", h.GetProductByID)
	}
}
