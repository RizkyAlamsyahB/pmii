package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// MockDashboardRepository adalah mock untuk DashboardRepository
type MockDashboardRepository struct {
	GetUniqueVisitorsFunc       func(startDate, endDate time.Time) (int64, error)
	GetPageViewsFunc            func(startDate, endDate time.Time) (int64, error)
	GetTotalPublishedPostsFunc  func() (int64, error)
	GetTopCategoryFunc          func(startDate, endDate time.Time) (*repository.CategoryStats, error)
	GetCategoryDistributionFunc func(startDate, endDate time.Time) ([]repository.CategoryStats, error)
	GetDailyVisitorsFunc        func(startDate, endDate time.Time) ([]repository.DailyVisitors, error)
	GetTopArticlesFunc          func(startDate, endDate time.Time, limit int) ([]repository.ArticleStats, error)
	GetAvailablePeriodsFunc     func() ([]time.Time, error)
	GetRecentActivityLogsFunc   func(limit int) ([]repository.ActivityLogStats, error)
}

func (m *MockDashboardRepository) GetUniqueVisitors(startDate, endDate time.Time) (int64, error) {
	if m.GetUniqueVisitorsFunc != nil {
		return m.GetUniqueVisitorsFunc(startDate, endDate)
	}
	return 0, nil
}

func (m *MockDashboardRepository) GetPageViews(startDate, endDate time.Time) (int64, error) {
	if m.GetPageViewsFunc != nil {
		return m.GetPageViewsFunc(startDate, endDate)
	}
	return 0, nil
}

func (m *MockDashboardRepository) GetTotalPublishedPosts() (int64, error) {
	if m.GetTotalPublishedPostsFunc != nil {
		return m.GetTotalPublishedPostsFunc()
	}
	return 0, nil
}

func (m *MockDashboardRepository) GetTopCategory(startDate, endDate time.Time) (*repository.CategoryStats, error) {
	if m.GetTopCategoryFunc != nil {
		return m.GetTopCategoryFunc(startDate, endDate)
	}
	return nil, nil
}

func (m *MockDashboardRepository) GetCategoryDistribution(startDate, endDate time.Time) ([]repository.CategoryStats, error) {
	if m.GetCategoryDistributionFunc != nil {
		return m.GetCategoryDistributionFunc(startDate, endDate)
	}
	return nil, nil
}

func (m *MockDashboardRepository) GetDailyVisitors(startDate, endDate time.Time) ([]repository.DailyVisitors, error) {
	if m.GetDailyVisitorsFunc != nil {
		return m.GetDailyVisitorsFunc(startDate, endDate)
	}
	return nil, nil
}

func (m *MockDashboardRepository) GetTopArticles(startDate, endDate time.Time, limit int) ([]repository.ArticleStats, error) {
	if m.GetTopArticlesFunc != nil {
		return m.GetTopArticlesFunc(startDate, endDate, limit)
	}
	return nil, nil
}

func (m *MockDashboardRepository) GetAvailablePeriods() ([]time.Time, error) {
	if m.GetAvailablePeriodsFunc != nil {
		return m.GetAvailablePeriodsFunc()
	}
	return nil, nil
}

func (m *MockDashboardRepository) GetRecentActivityLogs(limit int) ([]repository.ActivityLogStats, error) {
	if m.GetRecentActivityLogsFunc != nil {
		return m.GetRecentActivityLogsFunc(limit)
	}
	return nil, nil
}

// TestGetDashboard_Success menguji GetDashboard dengan data lengkap
func TestGetDashboard_Success(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		GetUniqueVisitorsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 11000, nil
		},
		GetPageViewsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 21000, nil
		},
		GetTotalPublishedPostsFunc: func() (int64, error) {
			return 500, nil
		},
		GetTopCategoryFunc: func(startDate, endDate time.Time) (*repository.CategoryStats, error) {
			return &repository.CategoryStats{
				ID:    1,
				Name:  "News",
				Views: 120000,
			}, nil
		},
		GetCategoryDistributionFunc: func(startDate, endDate time.Time) ([]repository.CategoryStats, error) {
			return []repository.CategoryStats{
				{ID: 1, Name: "News", Views: 50000},
				{ID: 2, Name: "Opini", Views: 30000},
				{ID: 3, Name: "Islamic", Views: 20000},
			}, nil
		},
		GetDailyVisitorsFunc: func(startDate, endDate time.Time) ([]repository.DailyVisitors, error) {
			return []repository.DailyVisitors{
				{Date: time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC), Visitors: 350},
				{Date: time.Date(2025, 12, 2, 0, 0, 0, 0, time.UTC), Visitors: 380},
			}, nil
		},
		GetTopArticlesFunc: func(startDate, endDate time.Time, limit int) ([]repository.ArticleStats, error) {
			return []repository.ArticleStats{
				{ID: 1, Title: "Santri Future Forum", Views: 3509},
				{ID: 2, Title: "GusDur Pahlawan", Views: 2800},
			}, nil
		},
	}

	dashboardService := NewDashboardService(mockRepo)
	result, err := dashboardService.GetDashboard(context.Background(), 2025, 12, false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.Period != "Desember 2025" {
		t.Errorf("Expected period 'Desember 2025', got '%s'", result.Period)
	}

	if result.Summary.UniqueVisitors != 11000 {
		t.Errorf("Expected unique visitors 11000, got %d", result.Summary.UniqueVisitors)
	}

	if result.Summary.PageViews != 21000 {
		t.Errorf("Expected page views 21000, got %d", result.Summary.PageViews)
	}

	if result.Summary.TotalPosts != 500 {
		t.Errorf("Expected total posts 500, got %d", result.Summary.TotalPosts)
	}

	if result.TopCategory == nil || result.TopCategory.Name != "News" {
		t.Errorf("Expected top category 'News', got %v", result.TopCategory)
	}

	if len(result.CategoryDistribution) != 3 {
		t.Errorf("Expected 3 categories, got %d", len(result.CategoryDistribution))
	}

	if len(result.VisitorsTrend) != 2 {
		t.Errorf("Expected 2 trend points, got %d", len(result.VisitorsTrend))
	}

	if len(result.TopArticles) != 2 {
		t.Errorf("Expected 2 top articles, got %d", len(result.TopArticles))
	}
}

// TestGetDashboard_EmptyData menguji GetDashboard tanpa data
func TestGetDashboard_EmptyData(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		GetUniqueVisitorsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 0, nil
		},
		GetPageViewsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 0, nil
		},
		GetTotalPublishedPostsFunc: func() (int64, error) {
			return 0, nil
		},
		GetTopCategoryFunc: func(startDate, endDate time.Time) (*repository.CategoryStats, error) {
			return nil, nil // No data
		},
		GetCategoryDistributionFunc: func(startDate, endDate time.Time) ([]repository.CategoryStats, error) {
			return []repository.CategoryStats{}, nil
		},
		GetDailyVisitorsFunc: func(startDate, endDate time.Time) ([]repository.DailyVisitors, error) {
			return []repository.DailyVisitors{}, nil
		},
		GetTopArticlesFunc: func(startDate, endDate time.Time, limit int) ([]repository.ArticleStats, error) {
			return []repository.ArticleStats{}, nil
		},
	}

	dashboardService := NewDashboardService(mockRepo)
	result, err := dashboardService.GetDashboard(context.Background(), 2025, 12, false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.TopCategory != nil {
		t.Errorf("Expected nil top category, got %v", result.TopCategory)
	}
}

// TestGetDashboard_RepositoryError menguji GetDashboard saat repository error
func TestGetDashboard_RepositoryError(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		GetUniqueVisitorsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 0, errors.New("database error")
		},
	}

	dashboardService := NewDashboardService(mockRepo)
	_, err := dashboardService.GetDashboard(context.Background(), 2025, 12, false)

	if err == nil {
		t.Error("Expected error, got nil")
	}
}

// TestGetDashboard_AdminWithActivityLogs menguji dashboard untuk admin dengan activity logs
func TestGetDashboard_AdminWithActivityLogs(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		GetUniqueVisitorsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 1000, nil
		},
		GetPageViewsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 5000, nil
		},
		GetTotalPublishedPostsFunc: func() (int64, error) {
			return 50, nil
		},
		GetTopCategoryFunc: func(startDate, endDate time.Time) (*repository.CategoryStats, error) {
			return nil, nil
		},
		GetCategoryDistributionFunc: func(startDate, endDate time.Time) ([]repository.CategoryStats, error) {
			return []repository.CategoryStats{}, nil
		},
		GetDailyVisitorsFunc: func(startDate, endDate time.Time) ([]repository.DailyVisitors, error) {
			return []repository.DailyVisitors{}, nil
		},
		GetTopArticlesFunc: func(startDate, endDate time.Time, limit int) ([]repository.ArticleStats, error) {
			return []repository.ArticleStats{}, nil
		},
		GetRecentActivityLogsFunc: func(limit int) ([]repository.ActivityLogStats, error) {
			return []repository.ActivityLogStats{
				{
					ID:         1,
					UserName:   "Admin User",
					ActionType: "CREATE",
					Module:     "posts",
					CreatedAt:  time.Date(2025, 6, 15, 10, 30, 0, 0, time.UTC),
				},
				{
					ID:         2,
					UserName:   "Admin User",
					ActionType: "UPDATE",
					Module:     "categories",
					CreatedAt:  time.Date(2025, 6, 15, 9, 0, 0, 0, time.UTC),
				},
			}, nil
		},
	}

	dashboardService := NewDashboardService(mockRepo)
	result, err := dashboardService.GetDashboard(context.Background(), 2025, 6, true)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Admin should have activity logs
	if len(result.ActivityLogs) != 2 {
		t.Errorf("Expected 2 activity logs for admin, got %d", len(result.ActivityLogs))
	}

	if result.ActivityLogs[0].ActionType != "CREATE" {
		t.Errorf("Expected first activity action 'CREATE', got '%s'", result.ActivityLogs[0].ActionType)
	}
}

// TestGetDashboard_AuthorWithoutActivityLogs menguji dashboard untuk author tanpa activity logs
func TestGetDashboard_AuthorWithoutActivityLogs(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		GetUniqueVisitorsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 1000, nil
		},
		GetPageViewsFunc: func(startDate, endDate time.Time) (int64, error) {
			return 5000, nil
		},
		GetTotalPublishedPostsFunc: func() (int64, error) {
			return 50, nil
		},
		GetTopCategoryFunc: func(startDate, endDate time.Time) (*repository.CategoryStats, error) {
			return nil, nil
		},
		GetCategoryDistributionFunc: func(startDate, endDate time.Time) ([]repository.CategoryStats, error) {
			return []repository.CategoryStats{}, nil
		},
		GetDailyVisitorsFunc: func(startDate, endDate time.Time) ([]repository.DailyVisitors, error) {
			return []repository.DailyVisitors{}, nil
		},
		GetTopArticlesFunc: func(startDate, endDate time.Time, limit int) ([]repository.ArticleStats, error) {
			return []repository.ArticleStats{}, nil
		},
		GetRecentActivityLogsFunc: func(limit int) ([]repository.ActivityLogStats, error) {
			// This should not be called for author
			t.Error("GetRecentActivityLogs should not be called for non-admin users")
			return nil, nil
		},
	}

	dashboardService := NewDashboardService(mockRepo)
	result, err := dashboardService.GetDashboard(context.Background(), 2025, 6, false)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	// Author should NOT have activity logs
	if len(result.ActivityLogs) != 0 {
		t.Errorf("Expected 0 activity logs for author, got %d", len(result.ActivityLogs))
	}
}

// TestCalculatePercentageChange menguji perhitungan persentase perubahan
func TestCalculatePercentageChange(t *testing.T) {
	tests := []struct {
		name     string
		previous int64
		current  int64
		expected float64
	}{
		{"Increase 100%", 100, 200, 100.0},
		{"Decrease 50%", 200, 100, -50.0},
		{"No change", 100, 100, 0.0},
		{"From zero", 0, 100, 100.0},
		{"To zero", 100, 0, -100.0},
		{"Both zero", 0, 0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculatePercentageChange(tt.previous, tt.current)
			if result != tt.expected {
				t.Errorf("Expected %f, got %f", tt.expected, result)
			}
		})
	}
}

// TestGetAvailablePeriods_Success menguji GetAvailablePeriods
func TestGetAvailablePeriods_Success(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		GetAvailablePeriodsFunc: func() ([]time.Time, error) {
			return []time.Time{
				time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
				time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
			}, nil
		},
	}

	dashboardService := NewDashboardService(mockRepo)
	result, err := dashboardService.GetAvailablePeriods(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(result.Periods) != 2 {
		t.Errorf("Expected 2 periods, got %d", len(result.Periods))
	}

	if result.Periods[0].Label != "Desember 2025" {
		t.Errorf("Expected 'Desember 2025', got '%s'", result.Periods[0].Label)
	}
}

// TestGetAvailablePeriods_Empty menguji GetAvailablePeriods tanpa data
func TestGetAvailablePeriods_Empty(t *testing.T) {
	mockRepo := &MockDashboardRepository{
		GetAvailablePeriodsFunc: func() ([]time.Time, error) {
			return []time.Time{}, nil
		},
	}

	dashboardService := NewDashboardService(mockRepo)
	result, err := dashboardService.GetAvailablePeriods(context.Background())

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Harus ada minimal 1 periode (bulan saat ini sebagai default)
	if len(result.Periods) < 1 {
		t.Errorf("Expected at least 1 period, got %d", len(result.Periods))
	}
}
