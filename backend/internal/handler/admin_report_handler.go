package handler

import (
	"awsome-shop/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AdminReportHandler handles admin report requests
type AdminReportHandler struct {
	pointsService     *service.PointsService
	redemptionService *service.RedemptionService
}

// NewAdminReportHandler creates a new AdminReportHandler instance
func NewAdminReportHandler(pointsService *service.PointsService, redemptionService *service.RedemptionService) *AdminReportHandler {
	return &AdminReportHandler{
		pointsService:     pointsService,
		redemptionService: redemptionService,
	}
}

// GetPointsGrantsReport gets the points grants report
// GET /api/v1/admin/reports/points-grants
func (h *AdminReportHandler) GetPointsGrantsReport(c *gin.Context) {
	stats, err := h.pointsService.GetGrantTransactionsReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve points grants report",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"grants": stats,
	})
}

// GetPointsBalancesReport gets the points balances report
// GET /api/v1/admin/reports/points-balances
func (h *AdminReportHandler) GetPointsBalancesReport(c *gin.Context) {
	stats, err := h.pointsService.GetPointsBalancesReport()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve points balances report",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"balances": stats,
	})
}

// GetRedemptionsReport gets the redemptions report
// GET /api/v1/admin/reports/redemptions
func (h *AdminReportHandler) GetRedemptionsReport(c *gin.Context) {
	stats, err := h.redemptionService.GetRedemptionStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve redemptions report",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"redemptions": stats,
	})
}

// RegisterRoutes registers admin report routes
func (h *AdminReportHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware, adminMiddleware gin.HandlerFunc) {
	admin := router.Group("/admin/reports")
	admin.Use(authMiddleware, adminMiddleware)
	{
		admin.GET("/points-grants", h.GetPointsGrantsReport)
		admin.GET("/points-balances", h.GetPointsBalancesReport)
		admin.GET("/redemptions", h.GetRedemptionsReport)
	}
}
