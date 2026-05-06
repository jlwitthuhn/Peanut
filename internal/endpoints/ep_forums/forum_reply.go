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

func registerForumReplyHandlers(mux *http.ServeMux, forumsService service.ForumService, forumThreadService service.ForumThreadService) {
	getReplyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		threadId := r.PathValue("threadId")

		thread, err := forumThreadService.GetThreadRowById(r.Context(), threadId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load thread.", w, r)
			return
		}
		if thread == nil {
			ep_util.RenderErrorHttp404NotFoundWithMessage("Page not found.", w, r)
			return
		}
		if thread.Visibility != "Public" {
			ep_util.RenderErrorHttp404NotFoundWithMessage("Page not found.", w, r)
			return
		}

		forum, err := forumsService.GetForumRowById(r.Context(), thread.ForumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}

		readable, err := forumsService.IsForumReadable(r.Context(), thread.ForumId)
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

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["ThreadId"] = thread.Id
		templateCtx["ThreadTitle"] = thread.Title
		templateCtx["ForumId"] = forum.Id
		templateCtx["ForumName"] = forum.Name
		templateCtx["SectionId"] = section.Id
		templateCtx["SectionName"] = section.Name
		ep_util.RenderTemplate("view_forum/reply", templateCtx, w, r)
	})
	mux.Handle("GET /forum/thread/{threadId}/reply", getReplyHandler)
}
