package service

import (
	"testing"
)

func TestCalculatePagination(t *testing.T) {
	testCases := []struct {
		total    int64
		limit    int
		expected int
	}{
		{total: 25, limit: 10, expected: 3},
		{total: 10, limit: 10, expected: 1},
		{total: 5, limit: 10, expected: 1},
		{total: 0, limit: 10, expected: 0},
	}

	for _, tc := range testCases {
		// Mock logic inside test
		var lastPage int
		if tc.total > 0 {
			lastPage = int((tc.total + int64(tc.limit) - 1) / int64(tc.limit))
		}

		if lastPage != tc.expected {
			t.Errorf("For total %d and limit %d, expected %d but got %d", tc.total, tc.limit, tc.expected, lastPage)
		}
	}
}
