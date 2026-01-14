package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/service"
	"github.com/gin-gonic/gin"
)

type ActivityLogHandler struct {
	svc service.ActivityLogService
}

func NewActivityLogHandler(svc service.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{svc: svc}
}

// GetActivityLogs handles GET /v1/activity-logs
// Query Params:
//   - page: int (default: 1)
//   - limit: int (default: 30)
//   - user_id: int (optional) - filter by user
//   - module: string (optional) - filter by module (user, post, category, etc.)
//   - action_type: string (optional) - filter by action type (create, update, delete, etc.)
//   - start_date: string (optional) - filter from date (format: 2006-01-02)
//   - end_date: string (optional) - filter to date (format: 2006-01-02)
//   - search: string (optional) - search in description
func (h *ActivityLogHandler) GetActivityLogs(c *gin.Context) {
	// Parse pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "30"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 30
	}

	// Build filter params
	filter := service.ActivityLogFilterParams{
		Search: c.Query("search"),
	}

	// Parse user_id filter
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.Atoi(userIDStr); err == nil {
			filter.UserID = &userID
		}
	}

	// Parse module filter
	if module := c.Query("module"); module != "" {
		filter.Module = &module
	}

	// Parse action_type filter
	if actionType := c.Query("action_type"); actionType != "" {
		filter.ActionType = &actionType
	}

	// Parse date range filters
	if startDateStr := c.Query("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			filter.StartDate = &startDate
		}
	}

	if endDateStr := c.Query("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set to end of day
			endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)
			filter.EndDate = &endDate
		}
	}

	// Call service
	logs, lastPage, total, err := h.svc.GetActivityLogs(page, limit, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, responses.ErrorResponse(500, "Gagal mengambil data activity logs"))
		return
	}

	c.JSON(http.StatusOK, responses.SuccessResponseWithPagination(
		200, "Activity logs berhasil dimuat", logs, page, limit, total, lastPage,
	))
}
