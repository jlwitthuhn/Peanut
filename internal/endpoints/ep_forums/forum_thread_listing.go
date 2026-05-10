// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"fmt"
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/service"
	"sort"
	"strconv"
)

type paginationItem struct {
	Number     int
	Url        string
	IsCurrent  bool
	IsEllipsis bool
}

type paginationWidget struct {
	Show    bool
	Items   []paginationItem
	HasPrev bool
	PrevUrl string
	HasNext bool
	NextUrl string
}

func buildThreadListingPagination(forumId string, currentPage int, totalPages int) paginationWidget {
	pageUrl := func(p int) string {
		if p == 1 {
			return "/forum/index/" + forumId
		}
		return fmt.Sprintf("/forum/index/%s/page/%d", forumId, p)
	}

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

	items := make([]paginationItem, 0, len(sorted))
	prev := 0
	for _, p := range sorted {
		if prev != 0 && p > prev+1 {
			items = append(items, paginationItem{IsEllipsis: true})
		}
		items = append(items, paginationItem{
			Number:    p,
			Url:       pageUrl(p),
			IsCurrent: p == currentPage,
		})
		prev = p
	}

	widget := paginationWidget{
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

func registerForumThreadListingHandlers(mux *http.ServeMux, forumsService service.ForumService, forumThreadService service.ForumThreadService) {
	renderPage := func(w http.ResponseWriter, r *http.Request, forumId string, pageNum int) {
		forum, err := forumsService.GetForumRowById(r.Context(), forumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}
		if forum == nil {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		readable, err := forumsService.IsForumReadable(r.Context(), forumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}
		if !readable {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		section, err := forumsService.GetSectionRowById(r.Context(), forum.SectionId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}

		page, err := forumThreadService.GetForumThreadListingPublic(r.Context(), forumId, pageNum)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load threads.", w, r)
			return
		}
		if page == nil {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		writable, err := forumThreadService.CanPostThreadInForum(r.Context(), forumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["ForumId"] = forum.Id
		templateCtx["ForumName"] = forum.Name
		templateCtx["SectionId"] = section.Id
		templateCtx["SectionName"] = section.Name
		templateCtx["Threads"] = page.Threads
		templateCtx["CanPostThread"] = writable
		templateCtx["Pagination"] = buildThreadListingPagination(forum.Id, page.CurrentPage, page.TotalPages)
		ep_util.RenderTemplate("view_forum/thread_listing", templateCtx, w, r)
	}

	getForumIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forumId := r.PathValue("forumId")
		renderPage(w, r, forumId, 1)
	})
	mux.Handle("GET /forum/index/{forumId}", getForumIndexHandler)

	getForumIndexPageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forumId := r.PathValue("forumId")
		pageStr := r.PathValue("pageNum")
		pageNum, err := strconv.Atoi(pageStr)
		if err != nil || pageNum < 1 {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}
		if pageNum == 1 {
			http.Redirect(w, r, "/forum/index/"+forumId, http.StatusMovedPermanently)
			return
		}
		renderPage(w, r, forumId, pageNum)
	})
	mux.Handle("GET /forum/index/{forumId}/page/{pageNum}", getForumIndexPageHandler)
}
