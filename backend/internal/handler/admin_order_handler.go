package handler

import (
	"awsome-shop/internal/service"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminOrderHandler handles admin order management requests
type AdminOrderHandler struct {
	redemptionService *service.RedemptionService
}

// NewAdminOrderHandler creates a new AdminOrderHandler instance
func NewAdminOrderHandler(redemptionService *service.RedemptionService) *AdminOrderHandler {
	return &AdminOrderHandler{
		redemptionService: redemptionService,
	}
}

// GetAllOrders gets all redemption orders
// GET /api/v1/admin/orders
func (h *AdminOrderHandler) GetAllOrders(c *gin.Context) {
	// Get optional filters
	var status *string
	var userID *uint

	if statusStr := c.Query("status"); statusStr != "" {
		status = &statusStr
	}

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		var uid uint
		if _, err := fmt.Sscanf(userIDStr, "%d", &uid); err == nil {
			userID = &uid
		}
	}

	orders, err := h.redemptionService.ListOrders(status, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve orders",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}

// BatchUpdateOrderStatusRequest represents a request to batch update order status
type BatchUpdateOrderStatusRequest struct {
	OrderNumbers string `json:"order_numbers" binding:"required"`
	Status       string `json:"status" binding:"required,oneof=preparing delivered"`
}

// BatchUpdateOrderStatus updates multiple orders' status
// PUT /api/v1/admin/orders/batch-status
func (h *AdminOrderHandler) BatchUpdateOrderStatus(c *gin.Context) {
	var req BatchUpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format. Status must be 'preparing' or 'delivered'",
		})
		return
	}

	// Parse order numbers
	orderNumbers := h.redemptionService.ParseBatchOrderNumbers(req.OrderNumbers)
	if len(orderNumbers) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No valid order numbers provided",
		})
		return
	}

	err := h.redemptionService.BatchUpdateOrderStatus(orderNumbers, req.Status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"count":   len(orderNumbers),
	})
}

// RegisterRoutes registers admin order routes
func (h *AdminOrderHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware, adminMiddleware gin.HandlerFunc) {
	admin := router.Group("/admin/orders")
	admin.Use(authMiddleware, adminMiddleware)
	{
		admin.GET("", h.GetAllOrders)
		admin.PUT("/batch-status", h.BatchUpdateOrderStatus)
	}
}
