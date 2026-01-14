package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

// DashboardHandler handles HTTP requests untuk dashboard endpoints
type DashboardHandler struct {
	dashboardService service.DashboardService
}

// NewDashboardHandler constructor untuk DashboardHandler
func NewDashboardHandler(dashboardService service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

// GetDashboard handles GET /dashboard
// Menampilkan data analytics dashboard dengan filter periode
// Query params: year (int), month (int) - default: bulan saat ini
// Admin mendapatkan tambahan activity_logs
func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	// Parse query parameters
	now := time.Now()
	year := now.Year()
	month := int(now.Month())

	// Parse year dari query param
	if yearStr := c.Query("year"); yearStr != "" {
		if parsedYear, err := strconv.Atoi(yearStr); err == nil {
			year = parsedYear
		}
	}

	// Parse month dari query param
	if monthStr := c.Query("month"); monthStr != "" {
		if parsedMonth, err := strconv.Atoi(monthStr); err == nil {
			month = parsedMonth
		}
	}

	// Check if user is admin (role = "1")
	isAdmin := false
	if userRole, exists := c.Get("user_role"); exists {
		if role, ok := userRole.(string); ok && role == "1" {
			isAdmin = true
		}
	}

	// Get dashboard data dari service
	dashboard, err := h.dashboardService.GetDashboard(c.Request.Context(), year, month, isAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data dashboard: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data dashboard berhasil diambil", dashboard))
}

// GetAvailablePeriods handles GET /admin/dashboard/periods
// Menampilkan list periode yang tersedia untuk filter dropdown
func (h *DashboardHandler) GetAvailablePeriods(c *gin.Context) {
	periods, err := h.dashboardService.GetAvailablePeriods(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data periode: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponse(200, "Data periode berhasil diambil", periods))
}
