// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_util

import "sort"

type PaginationItem struct {
	Number     int
	Url        string
	IsCurrent  bool
	IsEllipsis bool
}

type PaginationWidget struct {
	Show    bool
	Items   []PaginationItem
	HasPrev bool
	PrevUrl string
	HasNext bool
	NextUrl string
}

func BuildPagination(currentPage int, totalPages int, pageUrl func(int) string) PaginationWidget {
	pageSet := map[int]bool{}
	addPage := func(p int) {
		if p >= 1 && p <= totalPages {
			pageSet[p] = true
		}
	}
	addPage(1)
	addPage(2)
	addPage(currentPage - 2)
	addPage(currentPage - 1)
	addPage(currentPage)
	addPage(currentPage + 1)
	addPage(currentPage + 2)
	addPage(totalPages - 1)
	addPage(totalPages)

	sorted := make([]int, 0, len(pageSet))
	for p := range pageSet {
		sorted = append(sorted, p)
	}
	sort.Ints(sorted)

	items := make([]PaginationItem, 0, len(sorted))
	prev := 0
	for _, p := range sorted {
		if prev != 0 && p > prev+1 {
			items = append(items, PaginationItem{IsEllipsis: true})
		}
		items = append(items, PaginationItem{
			Number:    p,
			Url:       pageUrl(p),
			IsCurrent: p == currentPage,
		})
		prev = p
	}

	widget := PaginationWidget{
		Show:    totalPages > 1,
		Items:   items,
		HasPrev: currentPage > 1,
		HasNext: currentPage < totalPages,
	}
	if widget.HasPrev {
		widget.PrevUrl = pageUrl(currentPage - 1)
	}
	if widget.HasNext {
		widget.NextUrl = pageUrl(currentPage + 1)
	}
	return widget
}
