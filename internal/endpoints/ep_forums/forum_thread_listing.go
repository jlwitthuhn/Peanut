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
	"strconv"
)

func registerForumThreadListingHandlers(mux *http.ServeMux, forumsService service.ForumService, forumThreadService service.ForumThreadService) {
	pageUrl := func(forumId string, p int) string {
		if p == 1 {
			return "/forum/index/" + forumId
		}
		return fmt.Sprintf("/forum/index/%s/page/%d", forumId, p)
	}

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

		pagination := ep_util.BuildPagination(page.CurrentPage, page.TotalPages, func(p int) string {
			return pageUrl(forum.Id, p)
		})

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["ForumId"] = forum.Id
		templateCtx["ForumName"] = forum.Name
		templateCtx["SectionId"] = section.Id
		templateCtx["SectionName"] = section.Name
		templateCtx["Threads"] = page.Threads
		templateCtx["CanPostThread"] = writable
		templateCtx["Pagination"] = pagination
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
			http.Redirect(w, r, pageUrl(forumId, 1), http.StatusMovedPermanently)
			return
		}
		renderPage(w, r, forumId, pageNum)
	})
	mux.Handle("GET /forum/index/{forumId}/page/{pageNum}", getForumIndexPageHandler)
}
