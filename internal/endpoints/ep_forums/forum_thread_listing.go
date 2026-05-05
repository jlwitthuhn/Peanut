// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/service"
)

func registerForumThreadListingHandlers(mux *http.ServeMux, forumsService service.ForumService, forumThreadService service.ForumThreadService) {
	getForumIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forumId := r.PathValue("forumId")

		forum, err := forumsService.GetForumRowById(r.Context(), forumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}
		if forum == nil {
			ep_util.RenderErrorHttp404NotFoundWithMessage("Page not found.", w, r)
			return
		}

		readable, err := forumsService.IsForumReadable(r.Context(), forumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}
		if !readable {
			ep_util.RenderErrorHttp404NotFoundWithMessage("Page not found.", w, r)
			return
		}

		section, err := forumsService.GetSectionRowById(r.Context(), forum.SectionId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}

		threads, err := forumThreadService.GetForumThreadSummaryViewRowByForumIdPublic(r.Context(), forumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load threads.", w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["ForumId"] = forum.Id
		templateCtx["ForumName"] = forum.Name
		templateCtx["SectionId"] = section.Id
		templateCtx["SectionName"] = section.Name
		templateCtx["Threads"] = threads
		ep_util.RenderTemplate("view_forum/thread_listing", templateCtx, w, r)
	})
	mux.Handle("GET /forum/index/{forumId}", getForumIndexHandler)
}
