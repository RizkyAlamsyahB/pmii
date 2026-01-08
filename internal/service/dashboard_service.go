package service

import (
	"context"
	"fmt"
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/dto/responses"
	"github.com/garuda-labs-1/pmii-be/internal/repository"
)

// Nama bulan dalam bahasa Indonesia
var monthNames = []string{
	"", "Januari", "Februari", "Maret", "April", "Mei", "Juni",
	"Juli", "Agustus", "September", "Oktober", "November", "Desember",
}

// DashboardService interface untuk business logic dashboard
type DashboardService interface {
	// GetDashboard mendapatkan semua data dashboard untuk periode tertentu
	// isAdmin: jika true, akan menyertakan activity logs
	GetDashboard(ctx context.Context, year, month int, isAdmin bool) (*responses.DashboardResponse, error)

	// GetAvailablePeriods mendapatkan list periode yang tersedia untuk filter
	GetAvailablePeriods(ctx context.Context) (*responses.DashboardPeriodsResponse, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

// NewDashboardService constructor untuk DashboardService
func NewDashboardService(dashboardRepo repository.DashboardRepository) DashboardService {
	return &dashboardService{dashboardRepo: dashboardRepo}
}

// GetDashboard mendapatkan semua data dashboard
func (s *dashboardService) GetDashboard(ctx context.Context, year, month int, isAdmin bool) (*responses.DashboardResponse, error) {
	// Validasi input
	if year < 2020 || year > 2100 {
		year = time.Now().Year()
	}
	if month < 1 || month > 12 {
		month = int(time.Now().Month())
	}

	// Hitung range tanggal untuk bulan yang dipilih
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0) // Bulan berikutnya

	// Hitung range tanggal untuk bulan sebelumnya (untuk perbandingan)
	prevStartDate := startDate.AddDate(0, -1, 0)
	prevEndDate := startDate

	// 1. Get Summary Stats
	summary, err := s.getSummaryStats(startDate, endDate, prevStartDate, prevEndDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil summary stats: %w", err)
	}

	// 2. Get Top Category
	topCategory, err := s.getTopCategory(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil top category: %w", err)
	}

	// 3. Get Visitors Trend
	visitorsTrend, err := s.getVisitorsTrend(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil visitors trend: %w", err)
	}

	// 4. Get Category Distribution
	categoryDist, err := s.getCategoryDistribution(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil category distribution: %w", err)
	}

	// 5. Get Top Articles
	topArticles, err := s.getTopArticles(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil top articles: %w", err)
	}

	// Format period label
	periodLabel := fmt.Sprintf("%s %d", monthNames[month], year)

	response := &responses.DashboardResponse{
		Period:               periodLabel,
		Summary:              *summary,
		TopCategory:          topCategory,
		VisitorsTrend:        visitorsTrend,
		CategoryDistribution: categoryDist,
		TopArticles:          topArticles,
	}

	// 6. Get Activity Logs (hanya untuk Admin)
	if isAdmin {
		activityLogs, err := s.getRecentActivityLogs()
		if err != nil {
			// Log error tapi tidak gagalkan request
			activityLogs = []responses.ActivityLogItem{}
		}
		response.ActivityLogs = activityLogs
	}

	return response, nil
}

// getSummaryStats mendapatkan statistik ringkasan dengan perbandingan bulan lalu
func (s *dashboardService) getSummaryStats(startDate, endDate, prevStartDate, prevEndDate time.Time) (*responses.DashboardSummary, error) {
	// Current month stats
	uniqueVisitors, err := s.dashboardRepo.GetUniqueVisitors(startDate, endDate)
	if err != nil {
		return nil, err
	}

	pageViews, err := s.dashboardRepo.GetPageViews(startDate, endDate)
	if err != nil {
		return nil, err
	}

	totalPosts, err := s.dashboardRepo.GetTotalPublishedPosts()
	if err != nil {
		return nil, err
	}

	// Previous month stats (untuk perbandingan)
	prevUniqueVisitors, err := s.dashboardRepo.GetUniqueVisitors(prevStartDate, prevEndDate)
	if err != nil {
		return nil, err
	}

	prevPageViews, err := s.dashboardRepo.GetPageViews(prevStartDate, prevEndDate)
	if err != nil {
		return nil, err
	}

	// Calculate percentage changes
	uniqueVisitorsChange := calculatePercentageChange(prevUniqueVisitors, uniqueVisitors)
	pageViewsChange := calculatePercentageChange(prevPageViews, pageViews)

	return &responses.DashboardSummary{
		UniqueVisitors:       uniqueVisitors,
		UniqueVisitorsChange: uniqueVisitorsChange,
		PageViews:            pageViews,
		PageViewsChange:      pageViewsChange,
		TotalPosts:           totalPosts,
	}, nil
}

// getTopCategory mendapatkan kategori teratas
func (s *dashboardService) getTopCategory(startDate, endDate time.Time) (*responses.TopCategory, error) {
	stats, err := s.dashboardRepo.GetTopCategory(startDate, endDate)
	if err != nil {
		return nil, err
	}

	if stats == nil {
		return nil, nil
	}

	return &responses.TopCategory{
		Name:  stats.Name,
		Views: stats.Views,
	}, nil
}

// getVisitorsTrend mendapatkan trend visitor harian
func (s *dashboardService) getVisitorsTrend(startDate, endDate time.Time) ([]responses.VisitorsTrend, error) {
	dailyVisitors, err := s.dashboardRepo.GetDailyVisitors(startDate, endDate)
	if err != nil {
		return nil, err
	}

	result := make([]responses.VisitorsTrend, 0, len(dailyVisitors))
	for _, dv := range dailyVisitors {
		result = append(result, responses.VisitorsTrend{
			Date:     dv.Date.Format("2006-01-02"),
			Visitors: dv.Visitors,
		})
	}

	return result, nil
}

// getCategoryDistribution mendapatkan distribusi kategori
func (s *dashboardService) getCategoryDistribution(startDate, endDate time.Time) ([]responses.CategoryDistribution, error) {
	stats, err := s.dashboardRepo.GetCategoryDistribution(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Calculate total views
	var totalViews int64
	for _, cs := range stats {
		totalViews += cs.Views
	}

	result := make([]responses.CategoryDistribution, 0, len(stats))
	for _, cs := range stats {
		percentage := float64(0)
		if totalViews > 0 {
			percentage = float64(cs.Views) * 100.0 / float64(totalViews)
		}

		result = append(result, responses.CategoryDistribution{
			Name:       cs.Name,
			Views:      cs.Views,
			Percentage: roundToTwoDecimals(percentage),
		})
	}

	return result, nil
}

// getTopArticles mendapatkan top 5 artikel
func (s *dashboardService) getTopArticles(startDate, endDate time.Time) ([]responses.TopArticle, error) {
	stats, err := s.dashboardRepo.GetTopArticles(startDate, endDate, 5)
	if err != nil {
		return nil, err
	}

	result := make([]responses.TopArticle, 0, len(stats))
	for i, as := range stats {
		result = append(result, responses.TopArticle{
			Rank:  i + 1,
			ID:    as.ID,
			Title: as.Title,
			Views: as.Views,
		})
	}

	return result, nil
}

// GetAvailablePeriods mendapatkan list periode yang tersedia
func (s *dashboardService) GetAvailablePeriods(ctx context.Context) (*responses.DashboardPeriodsResponse, error) {
	periods, err := s.dashboardRepo.GetAvailablePeriods()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil periode: %w", err)
	}

	result := make([]responses.AvailablePeriod, 0, len(periods))

	// Jika tidak ada data, tambahkan bulan saat ini sebagai default
	if len(periods) == 0 {
		now := time.Now()
		result = append(result, responses.AvailablePeriod{
			Year:  now.Year(),
			Month: int(now.Month()),
			Label: fmt.Sprintf("%s %d", monthNames[int(now.Month())], now.Year()),
		})
	} else {
		for _, p := range periods {
			result = append(result, responses.AvailablePeriod{
				Year:  p.Year(),
				Month: int(p.Month()),
				Label: fmt.Sprintf("%s %d", monthNames[int(p.Month())], p.Year()),
			})
		}
	}

	return &responses.DashboardPeriodsResponse{
		Periods: result,
	}, nil
}

// getRecentActivityLogs mendapatkan activity logs terbaru (untuk admin)
func (s *dashboardService) getRecentActivityLogs() ([]responses.ActivityLogItem, error) {
	logs, err := s.dashboardRepo.GetRecentActivityLogs(5) // 5 activity terbaru
	if err != nil {
		return nil, err
	}

	// Mapping nama bulan Indonesia
	monthNamesShort := map[int]string{
		1: "Jan", 2: "Feb", 3: "Mar", 4: "Apr",
		5: "Mei", 6: "Jun", 7: "Jul", 8: "Agu",
		9: "Sep", 10: "Okt", 11: "Nov", 12: "Des",
	}

	result := make([]responses.ActivityLogItem, 0, len(logs))
	for _, log := range logs {
		// Format time: "10:24 • 03 Des 2025"
		monthShort := monthNamesShort[int(log.CreatedAt.Month())]
		timeFormatted := fmt.Sprintf("%s • %02d %s %d",
			log.CreatedAt.Format("15:04"),
			log.CreatedAt.Day(),
			monthShort,
			log.CreatedAt.Year(),
		)

		result = append(result, responses.ActivityLogItem{
			Name:  log.UserName,
			Title: log.Description,
			Time:  timeFormatted,
		})
	}

	return result, nil
}

// calculatePercentageChange menghitung persentase perubahan
func calculatePercentageChange(previous, current int64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100.0 // Naik 100% dari 0
		}
		return 0.0
	}

	change := float64(current-previous) * 100.0 / float64(previous)
	return roundToTwoDecimals(change)
}

// roundToTwoDecimals membulatkan ke 2 desimal
func roundToTwoDecimals(value float64) float64 {
	return float64(int(value*100)) / 100
}
