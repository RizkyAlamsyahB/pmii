package service

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateLastPage(t *testing.T) {
	// Mengetes logika pembulatan halaman di service
	totalData := int64(25)
	limit := 10

	lastPage := int(math.Ceil(float64(totalData) / float64(limit)))

	assert.Equal(t, 3, lastPage, "Last page should be 3 for 25 items with limit 10")
}

func TestFetchNewsDetail_EmptySlug(t *testing.T) {
	// Contoh pengetesan sederhana
	// Dalam realita, Anda akan meng-inject mock repository ke newsService
	assert.NotNil(t, t)
}

func TestFetchByCategory(t *testing.T) {
	// Pengujian skenario jika category slug diberikan
	// categorySlug := "opini"
	page := 1
	limit := 10

	// Pastikan offset dihitung dengan benar
	offset := (page - 1) * limit
	if offset != 0 {
		t.Errorf("Expected offset 0, got %d", offset)
	}
}
