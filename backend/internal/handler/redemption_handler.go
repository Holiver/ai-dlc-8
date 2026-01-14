package handler

import (
	"awsome-shop/internal/middleware"
	"awsome-shop/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// RedemptionHandler handles redemption related requests
type RedemptionHandler struct {
	redemptionService *service.RedemptionService
}

// NewRedemptionHandler creates a new RedemptionHandler instance
func NewRedemptionHandler(redemptionService *service.RedemptionService) *RedemptionHandler {
	return &RedemptionHandler{
		redemptionService: redemptionService,
	}
}

// RedeemProductRequest represents a request to redeem a product
type RedeemProductRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
}

// CreateRedemption creates a new redemption order
// POST /api/v1/redemptions
func (h *RedemptionHandler) CreateRedemption(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	var req RedeemProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Process redemption
	order, err := h.redemptionService.RedeemProduct(userID, req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order": order,
	})
}

// GetRedemptionHistory gets the current user's redemption history
// GET /api/v1/redemptions
func (h *RedemptionHandler) GetRedemptionHistory(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	orders, err := h.redemptionService.GetRedemptionHistory(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve redemption history",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

// GetRedemptionByID gets a redemption order by ID
// GET /api/v1/redemptions/:id
func (h *RedemptionHandler) GetRedemptionByID(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	var uri struct {
		ID uint `uri:"id" binding:"required"`
	}

	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order, err := h.redemptionService.GetOrderByID(uri.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	// Verify the order belongs to the current user
	if order.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

// RegisterRoutes registers redemption routes
func (h *RedemptionHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	redemptions := router.Group("/redemptions")
	redemptions.Use(authMiddleware)
	{
		redemptions.POST("", h.CreateRedemption)
		redemptions.GET("", h.GetRedemptionHistory)
		redemptions.GET("/:id", h.GetRedemptionByID)
	}
}
