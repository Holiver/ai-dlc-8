package handler

import (
	"awsome-shop/internal/middleware"
	"awsome-shop/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PointsHandler handles points related requests
type PointsHandler struct {
	pointsService *service.PointsService
}

// NewPointsHandler creates a new PointsHandler instance
func NewPointsHandler(pointsService *service.PointsService) *PointsHandler {
	return &PointsHandler{
		pointsService: pointsService,
	}
}

// GetPointsBalance gets the current user's points balance
// GET /api/v1/points/balance
func (h *PointsHandler) GetPointsBalance(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	balance, err := h.pointsService.GetPointsBalance(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve points balance",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balance": balance,
	})
}

// GetPointsTransactions gets the current user's points transaction history
// GET /api/v1/points/transactions
func (h *PointsHandler) GetPointsTransactions(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	// Get pagination parameters
	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	transactions, total, err := h.pointsService.GetPointsHistory(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve points transactions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"transactions": transactions,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
	})
}

// RegisterRoutes registers points routes
func (h *PointsHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	points := router.Group("/points")
	points.Use(authMiddleware)
	{
		points.GET("/balance", h.GetPointsBalance)
		points.GET("/transactions", h.GetPointsTransactions)
	}
}
