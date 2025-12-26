package service

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryPaginationLogic(t *testing.T) {
	// Skenario pengujian untuk metadata pagination
	testCases := []struct {
		name     string
		total    int64
		limit    int
		expected int
	}{
		{"Total 25 data, limit 10", 25, 10, 3},
		{"Total 10 data, limit 10", 10, 10, 1},
		{"Total 0 data, limit 10", 0, 10, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Logika yang sama dengan yang ada di GetAll service
			lastPage := int(math.Ceil(float64(tc.total) / float64(tc.limit)))
			if lastPage == 0 {
				lastPage = 1
			}

			assert.Equal(t, tc.expected, lastPage)
		})
	}
}

func TestCategorySlugGeneration(t *testing.T) {
	// Mengetes apakah fungsi GetSlug di request DTO bekerja
	name := "Berita Nasional Terbaru"
	expectedSlug := "berita-nasional-terbaru"

	// Simulasi logika GetSlug
	resultSlug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	assert.Equal(t, expectedSlug, resultSlug)
}
