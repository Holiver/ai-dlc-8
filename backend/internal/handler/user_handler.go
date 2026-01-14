package handler

import (
	"awsome-shop/internal/middleware"
	"awsome-shop/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user related requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetProfile gets the current user's profile
// GET /api/v1/users/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	user, err := h.userService.GetUserProfile(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

// UpdatePhoneRequest represents a request to update phone number
type UpdatePhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// UpdatePhone updates the current user's phone number
// PUT /api/v1/users/phone
func (h *UserHandler) UpdatePhone(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	var req UpdatePhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	err := h.userService.UpdatePhone(userID, req.Phone)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Phone number updated successfully",
	})
}

// RegisterRoutes registers user routes
func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	users := router.Group("/users")
	users.Use(authMiddleware)
	{
		users.GET("/profile", h.GetProfile)
		users.PUT("/phone", h.UpdatePhone)
	}
}
