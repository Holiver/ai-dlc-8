package handler

import (
	"awsome-shop/internal/middleware"
	"awsome-shop/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AdminUserHandler handles admin user management requests
type AdminUserHandler struct {
	userService *service.UserService
}

// NewAdminUserHandler creates a new AdminUserHandler instance
func NewAdminUserHandler(userService *service.UserService) *AdminUserHandler {
	return &AdminUserHandler{
		userService: userService,
	}
}

// CreateEmployeeRequest represents a request to create an employee
type CreateEmployeeRequest struct {
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
}

// CreateEmployee creates a new employee account
// POST /api/v1/admin/users
func (h *AdminUserHandler) CreateEmployee(c *gin.Context) {
	var req CreateEmployeeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	createReq := &service.CreateEmployeeRequest{
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
	}

	user, err := h.userService.CreateEmployee(createReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"user":    user,
		"message": "Employee created successfully. Initial password is the last 6 digits of phone number.",
	})
}

// SetEmployeeStatusRequest represents a request to set employee status
type SetEmployeeStatusRequest struct {
	IsActive bool `json:"is_active"`
}

// SetEmployeeStatus sets an employee's status (active/inactive for departure)
// PUT /api/v1/admin/users/:id/status
func (h *AdminUserHandler) SetEmployeeStatus(c *gin.Context) {
	operatorID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Operator ID not found in context",
		})
		return
	}

	// Get user ID from URL
	userIDStr := c.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req SetEmployeeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// If setting to inactive (departure), call SetEmployeeDeparture
	if !req.IsActive {
		err = h.userService.SetEmployeeDeparture(uint(userID), operatorID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Employee set to inactive. Points have been invalidated.",
		})
		return
	}

	// For reactivation, just update the status
	// Note: This is a simple implementation. In production, you might want more complex logic
	c.JSON(http.StatusBadRequest, gin.H{
		"error": "Reactivating employees is not supported in this version",
	})
}

// ListEmployees lists all employees
// GET /api/v1/admin/users
func (h *AdminUserHandler) ListEmployees(c *gin.Context) {
	// Get optional filter
	var isActive *bool
	if activeStr := c.Query("is_active"); activeStr != "" {
		if activeStr == "true" {
			val := true
			isActive = &val
		} else if activeStr == "false" {
			val := false
			isActive = &val
		}
	}

	users, err := h.userService.ListUsers(isActive)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}

// RegisterRoutes registers admin user routes
func (h *AdminUserHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware, adminMiddleware gin.HandlerFunc) {
	admin := router.Group("/admin/users")
	admin.Use(authMiddleware, adminMiddleware)
	{
		admin.POST("", h.CreateEmployee)
		admin.PUT("/:id/status", h.SetEmployeeStatus)
		admin.GET("", h.ListEmployees)
	}
}
