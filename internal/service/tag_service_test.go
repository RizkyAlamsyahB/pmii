package service

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTagPaginationCalculation(t *testing.T) {
	// Memastikan logika penentuan halaman terakhir benar
	testCases := []struct {
		total    int64
		limit    int
		expected int
	}{
		{25, 10, 3},
		{10, 10, 1},
		{0, 10, 1},
	}

	for _, tc := range testCases {
		lastPage := int(math.Ceil(float64(tc.total) / float64(tc.limit)))
		if lastPage == 0 {
			lastPage = 1
		}
		assert.Equal(t, tc.expected, lastPage)
	}
}
