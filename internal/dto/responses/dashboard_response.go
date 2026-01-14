package responses

// DashboardSummary ringkasan utama dashboard
type DashboardSummary struct {
	UniqueVisitors       int64   `json:"unique_visitors"`
	UniqueVisitorsChange float64 `json:"unique_visitors_change"` // persentase perubahan dari bulan lalu
	PageViews            int64   `json:"page_views"`
	PageViewsChange      float64 `json:"page_views_change"` // persentase perubahan dari bulan lalu
	TotalPosts           int64   `json:"total_posts"`
}

// TopCategory kategori dengan views terbanyak
type TopCategory struct {
	Name  string `json:"name"`
	Views int64  `json:"views"`
}

// VisitorsTrend data trend visitor per hari
type VisitorsTrend struct {
	Date     string `json:"date"`     // format: "2025-12-01"
	Visitors int64  `json:"visitors"` // unique visitors per hari
}

// CategoryDistribution distribusi views per kategori (untuk pie chart)
type CategoryDistribution struct {
	Name       string  `json:"name"`
	Views      int64   `json:"views"`
	Percentage float64 `json:"percentage"`
}

// TopArticle artikel dengan views terbanyak
type TopArticle struct {
	Rank  int    `json:"rank"`
	ID    int    `json:"id"`
	Title string `json:"title"`
	Views int64  `json:"views"`
}

// ActivityLogItem item activity log untuk dashboard admin
type ActivityLogItem struct {
	Name  string `json:"name"`  // Nama user
	Title string `json:"title"` // Deskripsi aktivitas
	Time  string `json:"time"`  // Format: "10:24 • 03 Des 2025"
}

// DashboardResponse response lengkap untuk dashboard
type DashboardResponse struct {
	Period               string                 `json:"period"` // "Desember 2025"
	Summary              DashboardSummary       `json:"summary"`
	TopCategory          *TopCategory           `json:"top_category"`
	VisitorsTrend        []VisitorsTrend        `json:"visitors_trend"`
	CategoryDistribution []CategoryDistribution `json:"category_distribution"`
	TopArticles          []TopArticle           `json:"top_articles"`
	ActivityLogs         []ActivityLogItem      `json:"activity_logs,omitempty"` // Hanya untuk Admin
}

// AvailablePeriod periode yang tersedia untuk filter
type AvailablePeriod struct {
	Year  int    `json:"year"`
	Month int    `json:"month"`
	Label string `json:"label"` // "Desember 2025"
}

// DashboardPeriodsResponse response untuk list periode yang tersedia
type DashboardPeriodsResponse struct {
	Periods []AvailablePeriod `json:"periods"`
}
