package repository

import (
	"time"

	"gorm.io/gorm"
)

// DashboardStats struktur untuk statistik dashboard
type DashboardStats struct {
	UniqueVisitors int64
	PageViews      int64
	TotalPosts     int64
}

// CategoryStats struktur untuk statistik per kategori
type CategoryStats struct {
	ID    int
	Name  string
	Views int64
}

// DailyVisitors struktur untuk visitor harian
type DailyVisitors struct {
	Date     time.Time
	Visitors int64
}

// ArticleStats struktur untuk statistik artikel
type ArticleStats struct {
	ID    int
	Title string
	Views int64
}

// ActivityLogStats struktur untuk activity log dashboard
type ActivityLogStats struct {
	ID          int
	UserName    string
	ActionType  string
	Module      string
	Description string
	CreatedAt   time.Time
}

// DashboardRepository interface untuk data layer dashboard
type DashboardRepository interface {
	// GetUniqueVisitors menghitung unique visitors berdasarkan IP dalam periode tertentu
	GetUniqueVisitors(startDate, endDate time.Time) (int64, error)

	// GetPageViews menghitung total page views dalam periode tertentu
	GetPageViews(startDate, endDate time.Time) (int64, error)

	// GetTotalPublishedPosts menghitung total post yang published
	GetTotalPublishedPosts() (int64, error)

	// GetTopCategory mendapatkan kategori dengan views terbanyak
	GetTopCategory(startDate, endDate time.Time) (*CategoryStats, error)

	// GetCategoryDistribution mendapatkan distribusi views per kategori
	GetCategoryDistribution(startDate, endDate time.Time) ([]CategoryStats, error)

	// GetDailyVisitors mendapatkan trend visitor harian
	GetDailyVisitors(startDate, endDate time.Time) ([]DailyVisitors, error)

	// GetTopArticles mendapatkan artikel dengan views terbanyak
	GetTopArticles(startDate, endDate time.Time, limit int) ([]ArticleStats, error)

	// GetAvailablePeriods mendapatkan periode yang tersedia (bulan-bulan yang ada datanya)
	GetAvailablePeriods() ([]time.Time, error)

	// GetRecentActivityLogs mendapatkan activity logs terbaru (untuk admin)
	GetRecentActivityLogs(limit int) ([]ActivityLogStats, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

// NewDashboardRepository constructor untuk DashboardRepository
func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetUniqueVisitors menghitung unique visitors berdasarkan IP
func (r *dashboardRepository) GetUniqueVisitors(startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("post_views").
		Where("viewed_at >= ? AND viewed_at < ?", startDate, endDate).
		Where("ip_address IS NOT NULL AND ip_address != ''").
		Distinct("ip_address").
		Count(&count).Error

	return count, err
}

// GetPageViews menghitung total page views
func (r *dashboardRepository) GetPageViews(startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Table("post_views").
		Where("viewed_at >= ? AND viewed_at < ?", startDate, endDate).
		Count(&count).Error

	return count, err
}

// GetTotalPublishedPosts menghitung total post yang published
func (r *dashboardRepository) GetTotalPublishedPosts() (int64, error) {
	var count int64
	err := r.db.Table("posts").
		Where("status = ?", "published").
		Where("deleted_at IS NULL").
		Count(&count).Error

	return count, err
}

// GetTopCategory mendapatkan kategori dengan views terbanyak
func (r *dashboardRepository) GetTopCategory(startDate, endDate time.Time) (*CategoryStats, error) {
	var result CategoryStats

	err := r.db.Table("post_views pv").
		Select("c.id, c.name, COUNT(pv.id) as views").
		Joins("JOIN posts p ON pv.post_id = p.id").
		Joins("JOIN categories c ON p.category_id = c.id").
		Where("pv.viewed_at >= ? AND pv.viewed_at < ?", startDate, endDate).
		Where("p.deleted_at IS NULL").
		Where("p.status = ?", "published").
		Group("c.id, c.name").
		Order("views DESC").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		return nil, err
	}

	// Jika tidak ada data
	if result.ID == 0 && result.Name == "" {
		return nil, nil
	}

	return &result, nil
}

// GetCategoryDistribution mendapatkan distribusi views per kategori
func (r *dashboardRepository) GetCategoryDistribution(startDate, endDate time.Time) ([]CategoryStats, error) {
	var results []CategoryStats

	err := r.db.Table("post_views pv").
		Select("c.id, c.name, COUNT(pv.id) as views").
		Joins("JOIN posts p ON pv.post_id = p.id").
		Joins("JOIN categories c ON p.category_id = c.id").
		Where("pv.viewed_at >= ? AND pv.viewed_at < ?", startDate, endDate).
		Where("p.deleted_at IS NULL").
		Where("p.status = ?", "published").
		Group("c.id, c.name").
		Order("views DESC").
		Scan(&results).Error

	return results, err
}

// GetDailyVisitors mendapatkan trend visitor harian
func (r *dashboardRepository) GetDailyVisitors(startDate, endDate time.Time) ([]DailyVisitors, error) {
	var results []DailyVisitors

	err := r.db.Table("post_views").
		Select("DATE(viewed_at) as date, COUNT(DISTINCT ip_address) as visitors").
		Where("viewed_at >= ? AND viewed_at < ?", startDate, endDate).
		Where("ip_address IS NOT NULL AND ip_address != ''").
		Group("DATE(viewed_at)").
		Order("date ASC").
		Scan(&results).Error

	return results, err
}

// GetTopArticles mendapatkan artikel dengan views terbanyak
func (r *dashboardRepository) GetTopArticles(startDate, endDate time.Time, limit int) ([]ArticleStats, error) {
	var results []ArticleStats

	err := r.db.Table("post_views pv").
		Select("p.id, p.title, COUNT(pv.id) as views").
		Joins("JOIN posts p ON pv.post_id = p.id").
		Where("pv.viewed_at >= ? AND pv.viewed_at < ?", startDate, endDate).
		Where("p.deleted_at IS NULL").
		Where("p.status = ?", "published").
		Group("p.id, p.title, p.slug").
		Order("views DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}

// GetAvailablePeriods mendapatkan periode yang tersedia
func (r *dashboardRepository) GetAvailablePeriods() ([]time.Time, error) {
	var results []time.Time

	err := r.db.Table("post_views").
		Select("DATE_TRUNC('month', viewed_at) as period").
		Group("DATE_TRUNC('month', viewed_at)").
		Order("period DESC").
		Limit(12). // Maksimal 12 bulan terakhir
		Pluck("period", &results).Error

	return results, err
}

// GetRecentActivityLogs mendapatkan activity logs terbaru
func (r *dashboardRepository) GetRecentActivityLogs(limit int) ([]ActivityLogStats, error) {
	var results []ActivityLogStats

	err := r.db.Table("activity_logs al").
		Select("al.id, u.full_name as user_name, al.action_type, al.module, al.description, al.created_at").
		Joins("JOIN users u ON al.user_id = u.id").
		Order("al.created_at DESC").
		Limit(limit).
		Scan(&results).Error

	return results, err
}
