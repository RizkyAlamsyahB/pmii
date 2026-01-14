package service

import (
	"math"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

type ActivityLogService interface {
	GetActivityLogs(page, limit int, filter ActivityLogFilterParams) ([]responses.ActivityLogResponse, int, int64, error)
}

// ActivityLogFilterParams is the service-level filter struct for activity logs
type ActivityLogFilterParams struct {
	UserID     *int
	Module     *string
	ActionType *string
	StartDate  *time.Time
	EndDate    *time.Time
	Search     string
}

type activityLogService struct {
	repo repository.ActivityLogRepository
}

func NewActivityLogService(repo repository.ActivityLogRepository) ActivityLogService {
	return &activityLogService{repo: repo}
}

func (s *activityLogService) GetActivityLogs(page, limit int, filter ActivityLogFilterParams) ([]responses.ActivityLogResponse, int, int64, error) {
	offset := (page - 1) * limit

	// Convert service filter to repository filter
	repoFilter := repository.ActivityLogFilter{
		UserID:     filter.UserID,
		Module:     filter.Module,
		ActionType: filter.ActionType,
		StartDate:  filter.StartDate,
		EndDate:    filter.EndDate,
		Search:     filter.Search,
	}

	logs, total, err := s.repo.GetActivityLogs(offset, limit, repoFilter)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := int(math.Ceil(float64(total) / float64(limit)))
	if lastPage < 1 {
		lastPage = 1
	}

	return logs, lastPage, total, nil
}
