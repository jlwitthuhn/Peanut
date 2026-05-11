// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_util

import (
	"fmt"
	"testing"
)

func pageUrl(p int) string {
	return fmt.Sprintf("/p/%d", p)
}

func TestBuildPagination_Visibility(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int
		totalPages  int
		wantShow    bool
	}{
		{"one page", 1, 1, false},
		{"two pages", 1, 2, true},
		{"many pages", 5, 10, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPagination(tt.currentPage, tt.totalPages, pageUrl)
			if got.Show != tt.wantShow {
				t.Errorf("Show: got %v, want %v", got.Show, tt.wantShow)
			}
		})
	}
}

func TestBuildPagination_PrevNext(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int
		totalPages  int
		wantHasPrev bool
		wantPrevUrl string
		wantHasNext bool
		wantNextUrl string
	}{
		{"first page of many", 1, 10, false, "", true, "/p/2"},
		{"middle page", 5, 10, true, "/p/4", true, "/p/6"},
		{"last page", 10, 10, true, "/p/9", false, ""},
		{"only page", 1, 1, false, "", false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPagination(tt.currentPage, tt.totalPages, pageUrl)
			if got.HasPrev != tt.wantHasPrev {
				t.Errorf("HasPrev: got %v, want %v", got.HasPrev, tt.wantHasPrev)
			}
			if got.PrevUrl != tt.wantPrevUrl {
				t.Errorf("PrevUrl: got %q, want %q", got.PrevUrl, tt.wantPrevUrl)
			}
			if got.HasNext != tt.wantHasNext {
				t.Errorf("HasNext: got %v, want %v", got.HasNext, tt.wantHasNext)
			}
			if got.NextUrl != tt.wantNextUrl {
				t.Errorf("NextUrl: got %q, want %q", got.NextUrl, tt.wantNextUrl)
			}
		})
	}
}

// itemsToString renders items compactly for easy diffing in test output.
// A normal page becomes "3", the current page becomes "[3]", and an ellipsis becomes "...".
func itemsToString(items []PaginationItem) string {
	s := ""
	for i, it := range items {
		if i > 0 {
			s += " "
		}
		switch {
		case it.IsEllipsis:
			s += "..."
		case it.IsCurrent:
			s += fmt.Sprintf("[%d]", it.Number)
		default:
			s += fmt.Sprintf("%d", it.Number)
		}
	}
	return s
}

func TestBuildPagination_ItemLayout(t *testing.T) {
	tests := []struct {
		name        string
		currentPage int
		totalPages  int
		want        string
	}{
		{"single page", 1, 1, "[1]"},
		{"few pages no ellipsis", 2, 4, "1 [2] 3 4"},
		{"five pages current middle", 3, 5, "1 2 [3] 4 5"},
		{"current at start with gap", 1, 20, "[1] 2 3 ... 19 20"},
		{"current at end with gap", 20, 20, "1 2 ... 18 19 [20]"},
		{"current in middle with two gaps", 10, 20, "1 2 ... 8 9 [10] 11 12 ... 19 20"},
		{"current near start no left gap", 4, 10, "1 2 3 [4] 5 6 ... 9 10"},
		{"current near end no right gap", 7, 10, "1 2 ... 5 6 [7] 8 9 10"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildPagination(tt.currentPage, tt.totalPages, pageUrl)
			rendered := itemsToString(got.Items)
			if rendered != tt.want {
				t.Errorf("\nlayout: %s\n   got: %s\n  want: %s", tt.name, rendered, tt.want)
			}
		})
	}
}

func TestBuildPagination_UrlMatchesNumber(t *testing.T) {
	got := BuildPagination(5, 10, pageUrl)
	for _, it := range got.Items {
		if it.IsEllipsis {
			if it.Url != "" {
				t.Errorf("ellipsis item should have empty Url, got %q", it.Url)
			}
			continue
		}
		want := fmt.Sprintf("/p/%d", it.Number)
		if it.Url != want {
			t.Errorf("page %d: Url got %q, want %q", it.Number, it.Url, want)
		}
	}
}
