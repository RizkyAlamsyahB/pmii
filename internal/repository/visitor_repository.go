package repository

import (
	"time"

	"github.com/garuda-labs-1/pmii-be/internal/domain"
	"gorm.io/gorm"
)

// VisitorRepository interface untuk visitor tracking
type VisitorRepository interface {
	// RecordVisit mencatat visitor baru (1 IP = 1 record per hari)
	RecordVisit(ipAddress string, platform string) error

	// HasVisitedToday mengecek apakah IP sudah tercatat hari ini
	HasVisitedToday(ipAddress string) (bool, error)

	// GetUniqueVisitors menghitung unique visitors dalam periode tertentu
	GetUniqueVisitors(startDate, endDate time.Time) (int64, error)

	// GetDailyVisitors mendapatkan jumlah visitor per hari
	GetDailyVisitors(startDate, endDate time.Time) ([]DailyVisitors, error)
}

type visitorRepository struct {
	db *gorm.DB
}

// NewVisitorRepository constructor untuk VisitorRepository
func NewVisitorRepository(db *gorm.DB) VisitorRepository {
	return &visitorRepository{db: db}
}

// RecordVisit mencatat visitor baru jika belum ada record hari ini
func (r *visitorRepository) RecordVisit(ipAddress string, platform string) error {
	// Cek apakah IP sudah tercatat hari ini
	hasVisited, err := r.HasVisitedToday(ipAddress)
	if err != nil {
		return err
	}

	// Jika sudah tercatat hari ini, skip insert
	if hasVisited {
		return nil
	}

	// Insert record baru
	visitor := domain.Visitor{
		IPAddress: ipAddress,
		Platform:  &platform,
		VisitedAt: time.Now(),
	}

	return r.db.Create(&visitor).Error
}

// HasVisitedToday mengecek apakah IP sudah tercatat hari ini
func (r *visitorRepository) HasVisitedToday(ipAddress string) (bool, error) {
	var count int64
	today := time.Now().Format("2006-01-02")

	err := r.db.Model(&domain.Visitor{}).
		Where("ip_address = ?", ipAddress).
		Where("DATE(visited_at) = ?", today).
		Count(&count).Error

	return count > 0, err
}

// GetUniqueVisitors menghitung total record dalam periode (sudah unique per hari)
func (r *visitorRepository) GetUniqueVisitors(startDate, endDate time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&domain.Visitor{}).
		Where("visited_at >= ? AND visited_at < ?", startDate, endDate).
		Count(&count).Error

	return count, err
}

// GetDailyVisitors mendapatkan jumlah visitor per hari
func (r *visitorRepository) GetDailyVisitors(startDate, endDate time.Time) ([]DailyVisitors, error) {
	var results []DailyVisitors

	err := r.db.Model(&domain.Visitor{}).
		Select("DATE(visited_at) as date, COUNT(*) as visitors").
		Where("visited_at >= ? AND visited_at < ?", startDate, endDate).
		Group("DATE(visited_at)").
		Order("date ASC").
		Scan(&results).Error

	return results, err
}
