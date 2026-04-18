package main

import "testing"

func TestClampPageWindow(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		page     int
		pageSize int
		total    int
		wantFrom int
		wantTo   int
	}{
		{name: "first page", page: 1, pageSize: 2, total: 5, wantFrom: 0, wantTo: 2},
		{name: "after last page", page: 10, pageSize: 2, total: 5, wantFrom: 5, wantTo: 5},
		{name: "very large page", page: int(^uint(0) >> 1), pageSize: 2, total: 5, wantFrom: 5, wantTo: 5},
		{name: "very large page size", page: 1, pageSize: int(^uint(0) >> 1), total: 5, wantFrom: 0, wantTo: 5},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			from, to := clampPageWindow(tt.page, tt.pageSize, tt.total)
			if from != tt.wantFrom || to != tt.wantTo {
				t.Fatalf("expected (%d,%d), got (%d,%d)", tt.wantFrom, tt.wantTo, from, to)
			}
		})
	}
}
