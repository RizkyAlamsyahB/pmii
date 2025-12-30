package service

import (
	"errors"
	"testing"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// MockActivityLogRepository adalah mock untuk ActivityLogRepository
type MockActivityLogRepository struct {
	GetActivityLogsFunc func(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error)
}

func (m *MockActivityLogRepository) GetActivityLogs(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
	if m.GetActivityLogsFunc != nil {
		return m.GetActivityLogsFunc(offset, limit, filter)
	}
	return nil, 0, errors.New("mock not configured")
}

// TestGetActivityLogs_Success menguji GetActivityLogs yang berhasil
func TestGetActivityLogs_Success(t *testing.T) {
	description := "Created new post"
	targetID := 1
	ipAddress := "127.0.0.1"
	userAgent := "Mozilla/5.0"

	mockLogs := []responses.ActivityLogResponse{
		{
			ID:          1,
			UserID:      1,
			User:        &responses.ActivityLogUserInfo{ID: 1, FullName: "Admin User", Email: "admin@example.com"},
			ActionType:  domain.ActionCreate,
			Module:      domain.ModulePost,
			Description: &description,
			TargetID:    &targetID,
			OldValue:    nil,
			NewValue:    map[string]any{"title": "New Post"},
			IPAddress:   &ipAddress,
			UserAgent:   &userAgent,
			CreatedAt:   time.Now(),
		},
		{
			ID:          2,
			UserID:      2,
			User:        &responses.ActivityLogUserInfo{ID: 2, FullName: "Author User", Email: "author@example.com"},
			ActionType:  domain.ActionUpdate,
			Module:      domain.ModulePost,
			Description: &description,
			TargetID:    &targetID,
			OldValue:    map[string]any{"title": "Old Title"},
			NewValue:    map[string]any{"title": "New Title"},
			IPAddress:   &ipAddress,
			UserAgent:   &userAgent,
			CreatedAt:   time.Now(),
		},
	}

	mockRepo := &MockActivityLogRepository{
		GetActivityLogsFunc: func(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
			if offset != 0 {
				t.Errorf("Expected offset 0, got %d", offset)
			}
			if limit != 10 {
				t.Errorf("Expected limit 10, got %d", limit)
			}
			return mockLogs, 2, nil
		},
	}

	service := NewActivityLogService(mockRepo)
	logs, lastPage, total, err := service.GetActivityLogs(1, 10, ActivityLogFilterParams{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(logs) != 2 {
		t.Errorf("Expected 2 logs, got %d", len(logs))
	}
	if lastPage != 1 {
		t.Errorf("Expected lastPage 1, got %d", lastPage)
	}
	if total != 2 {
		t.Errorf("Expected total 2, got %d", total)
	}
}

// TestGetActivityLogs_RepositoryError menguji GetActivityLogs ketika repository mengembalikan error
func TestGetActivityLogs_RepositoryError(t *testing.T) {
	mockRepo := &MockActivityLogRepository{
		GetActivityLogsFunc: func(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
			return nil, 0, errors.New("database connection failed")
		},
	}

	service := NewActivityLogService(mockRepo)
	logs, lastPage, total, err := service.GetActivityLogs(1, 10, ActivityLogFilterParams{})

	if err == nil {
		t.Error("Expected error, got nil")
	}
	if logs != nil {
		t.Errorf("Expected nil logs, got %v", logs)
	}
	if lastPage != 0 {
		t.Errorf("Expected lastPage 0, got %d", lastPage)
	}
	if total != 0 {
		t.Errorf("Expected total 0, got %d", total)
	}
}

// TestGetActivityLogs_EmptyResult menguji GetActivityLogs ketika tidak ada data
func TestGetActivityLogs_EmptyResult(t *testing.T) {
	mockRepo := &MockActivityLogRepository{
		GetActivityLogsFunc: func(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
			return []responses.ActivityLogResponse{}, 0, nil
		},
	}

	service := NewActivityLogService(mockRepo)
	logs, lastPage, total, err := service.GetActivityLogs(1, 10, ActivityLogFilterParams{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(logs) != 0 {
		t.Errorf("Expected 0 logs, got %d", len(logs))
	}
	if lastPage != 1 {
		t.Errorf("Expected lastPage 1 (minimum), got %d", lastPage)
	}
	if total != 0 {
		t.Errorf("Expected total 0, got %d", total)
	}
}

// TestGetActivityLogs_WithPagination menguji GetActivityLogs dengan pagination
func TestGetActivityLogs_WithPagination(t *testing.T) {
	mockRepo := &MockActivityLogRepository{
		GetActivityLogsFunc: func(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
			// Verifikasi offset dihitung dengan benar untuk page 3
			if offset != 20 {
				t.Errorf("Expected offset 20 for page 3 with limit 10, got %d", offset)
			}
			if limit != 10 {
				t.Errorf("Expected limit 10, got %d", limit)
			}
			return []responses.ActivityLogResponse{}, 50, nil
		},
	}

	service := NewActivityLogService(mockRepo)
	_, lastPage, total, err := service.GetActivityLogs(3, 10, ActivityLogFilterParams{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if lastPage != 5 {
		t.Errorf("Expected lastPage 5 (50 items / 10 per page), got %d", lastPage)
	}
	if total != 50 {
		t.Errorf("Expected total 50, got %d", total)
	}
}

// TestGetActivityLogs_LastPageCalculation menguji perhitungan lastPage dengan berbagai skenario
func TestGetActivityLogs_LastPageCalculation(t *testing.T) {
	testCases := []struct {
		name             string
		total            int64
		limit            int
		expectedLastPage int
	}{
		{"ExactDivision", 100, 10, 10},
		{"WithRemainder", 101, 10, 11},
		{"ZeroItems", 0, 10, 1},
		{"OneItem", 1, 10, 1},
		{"LimitEqualsTotal", 10, 10, 1},
		{"LimitGreaterThanTotal", 5, 10, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &MockActivityLogRepository{
				GetActivityLogsFunc: func(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
					return []responses.ActivityLogResponse{}, tc.total, nil
				},
			}

			service := NewActivityLogService(mockRepo)
			_, lastPage, _, err := service.GetActivityLogs(1, tc.limit, ActivityLogFilterParams{})

			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if lastPage != tc.expectedLastPage {
				t.Errorf("Expected lastPage %d, got %d", tc.expectedLastPage, lastPage)
			}
		})
	}
}

// TestGetActivityLogs_FilterConversion menguji konversi filter dari service ke repository
func TestGetActivityLogs_FilterConversion(t *testing.T) {
	userID := 1
	module := "post"
	actionType := "create"
	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()
	search := "test search"

	mockRepo := &MockActivityLogRepository{
		GetActivityLogsFunc: func(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
			// Verifikasi semua filter dikonversi dengan benar
			if filter.UserID == nil || *filter.UserID != userID {
				t.Errorf("Expected UserID %d, got %v", userID, filter.UserID)
			}
			if filter.Module == nil || *filter.Module != module {
				t.Errorf("Expected Module %s, got %v", module, filter.Module)
			}
			if filter.ActionType == nil || *filter.ActionType != actionType {
				t.Errorf("Expected ActionType %s, got %v", actionType, filter.ActionType)
			}
			if filter.StartDate == nil || !filter.StartDate.Equal(startDate) {
				t.Errorf("Expected StartDate %v, got %v", startDate, filter.StartDate)
			}
			if filter.EndDate == nil || !filter.EndDate.Equal(endDate) {
				t.Errorf("Expected EndDate %v, got %v", endDate, filter.EndDate)
			}
			if filter.Search != search {
				t.Errorf("Expected Search %s, got %s", search, filter.Search)
			}
			return []responses.ActivityLogResponse{}, 0, nil
		},
	}

	service := NewActivityLogService(mockRepo)
	serviceFilter := ActivityLogFilterParams{
		UserID:     &userID,
		Module:     &module,
		ActionType: &actionType,
		StartDate:  &startDate,
		EndDate:    &endDate,
		Search:     search,
	}

	_, _, _, err := service.GetActivityLogs(1, 10, serviceFilter)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestGetActivityLogs_NilFilters menguji GetActivityLogs dengan filter nil
func TestGetActivityLogs_NilFilters(t *testing.T) {
	mockRepo := &MockActivityLogRepository{
		GetActivityLogsFunc: func(offset, limit int, filter repository.ActivityLogFilter) ([]responses.ActivityLogResponse, int64, error) {
			// Verifikasi filter nil dilewatkan dengan benar
			if filter.UserID != nil {
				t.Error("Expected UserID to be nil")
			}
			if filter.Module != nil {
				t.Error("Expected Module to be nil")
			}
			if filter.ActionType != nil {
				t.Error("Expected ActionType to be nil")
			}
			if filter.StartDate != nil {
				t.Error("Expected StartDate to be nil")
			}
			if filter.EndDate != nil {
				t.Error("Expected EndDate to be nil")
			}
			if filter.Search != "" {
				t.Errorf("Expected Search to be empty, got %s", filter.Search)
			}
			return []responses.ActivityLogResponse{}, 0, nil
		},
	}

	service := NewActivityLogService(mockRepo)
	_, _, _, err := service.GetActivityLogs(1, 10, ActivityLogFilterParams{})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

// TestNewActivityLogService menguji pembuatan service baru
func TestNewActivityLogService(t *testing.T) {
	mockRepo := &MockActivityLogRepository{}
	service := NewActivityLogService(mockRepo)

	if service == nil {
		t.Error("Expected service to not be nil")
	}
}
