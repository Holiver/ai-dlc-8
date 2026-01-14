package handler

import (
	"awsome-shop/internal/middleware"
	"awsome-shop/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminPointsHandler handles admin points management requests
type AdminPointsHandler struct {
	pointsService *service.PointsService
}

// NewAdminPointsHandler creates a new AdminPointsHandler instance
func NewAdminPointsHandler(pointsService *service.PointsService) *AdminPointsHandler {
	return &AdminPointsHandler{
		pointsService: pointsService,
	}
}

// GrantPointsRequest represents a request to grant points
type GrantPointsRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Amount int    `json:"amount" binding:"required,min=1"`
	Reason string `json:"reason" binding:"required"`
}

// GrantPoints grants points to a user
// POST /api/v1/admin/points/grant
func (h *AdminPointsHandler) GrantPoints(c *gin.Context) {
	operatorID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Operator ID not found in context",
		})
		return
	}

	var req GrantPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	err := h.pointsService.GrantPoints(req.UserID, req.Amount, req.Reason, operatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Points granted successfully",
	})
}

// DeductPointsRequest represents a request to deduct points
type DeductPointsRequest struct {
	UserID uint   `json:"user_id" binding:"required"`
	Amount int    `json:"amount" binding:"required,min=1"`
	Reason string `json:"reason" binding:"required"`
}

// DeductPoints deducts points from a user
// POST /api/v1/admin/points/deduct
func (h *AdminPointsHandler) DeductPoints(c *gin.Context) {
	operatorID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Operator ID not found in context",
		})
		return
	}

	var req DeductPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	err := h.pointsService.DeductPoints(req.UserID, req.Amount, req.Reason, operatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Points deducted successfully",
	})
}

// BatchGrantPointsRequest represents a request to batch grant points
type BatchGrantPointsRequest struct {
	Markdown string `json:"markdown" binding:"required"`
}

// BatchGrantPoints grants points to multiple users
// POST /api/v1/admin/points/batch-grant
func (h *AdminPointsHandler) BatchGrantPoints(c *gin.Context) {
	operatorID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Operator ID not found in context",
		})
		return
	}

	var req BatchGrantPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	err := h.pointsService.BatchGrantPoints(req.Markdown, operatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Points granted to all users successfully",
	})
}

// RegisterRoutes registers admin points routes
func (h *AdminPointsHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware, adminMiddleware gin.HandlerFunc) {
	admin := router.Group("/admin/points")
	admin.Use(authMiddleware, adminMiddleware)
	{
		admin.POST("/grant", h.GrantPoints)
		admin.POST("/deduct", h.DeductPoints)
		admin.POST("/batch-grant", h.BatchGrantPoints)
	}
}
