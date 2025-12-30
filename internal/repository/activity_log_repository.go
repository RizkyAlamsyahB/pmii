package repository

import (
	"time"

	"github.com/garuda-labs-1/pmii-be/config"
	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
)

// ActivityLogFilter contains filter options for querying activity logs
type ActivityLogFilter struct {
	UserID     *int
	Module     *string
	ActionType *string
	StartDate  *time.Time
	EndDate    *time.Time
	Search     string
}

type ActivityLogRepository interface {
	GetActivityLogs(offset, limit int, filter ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error)
	Create(log *domain.ActivityLog) error
}

type activityLogRepository struct{}

func NewActivityLogRepository() ActivityLogRepository {
	return &activityLogRepository{}
}

func (r *activityLogRepository) GetActivityLogs(offset, limit int, filter ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
	var logs []domain.ActivityLog
	var total int64

	db := config.DB.Model(&domain.ActivityLog{}).Preload("User")

	// Apply filters
	if filter.UserID != nil {
		db = db.Where("user_id = ?", *filter.UserID)
	}
	if filter.Module != nil && *filter.Module != "" {
		db = db.Where("module = ?", *filter.Module)
	}
	if filter.ActionType != nil && *filter.ActionType != "" {
		db = db.Where("action_type = ?", *filter.ActionType)
	}
	if filter.StartDate != nil {
		db = db.Where("created_at >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		db = db.Where("created_at <= ?", *filter.EndDate)
	}
	if filter.Search != "" {
		db = db.Where("description ILIKE ?", "%"+filter.Search+"%")
	}

	// Count total before pagination
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination and ordering
	if err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	// Convert to response DTOs
	result := make([]responses.ActivityLogResponse, len(logs))
	for i, log := range logs {
		var userInfo *responses.ActivityLogUserInfo
		if log.User.ID != 0 {
			userInfo = &responses.ActivityLogUserInfo{
				ID:       log.User.ID,
				FullName: log.User.FullName,
				Email:    log.User.Email,
			}
		}

		result[i] = responses.ActivityLogResponse{
			ID:          log.ID,
			UserID:      log.UserID,
			User:        userInfo,
			ActionType:  log.ActionType,
			Module:      log.Module,
			Description: log.Description,
			TargetID:    log.TargetID,
			OldValue:    log.OldValue,
			NewValue:    log.NewValue,
			IPAddress:   log.IPAddress,
			UserAgent:   log.UserAgent,
			CreatedAt:   log.CreatedAt,
		}
	}

	return result, total, nil
}

// Create inserts a new activity log entry
func (r *activityLogRepository) Create(log *domain.ActivityLog) error {
	return config.DB.Create(log).Error
}
