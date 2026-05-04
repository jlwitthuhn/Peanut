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

func registerForumIndexHandlers(mux *http.ServeMux, forumsService service.ForumsService) {
	getForumIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forumId := r.PathValue("forumId")

		readable, err := forumsService.IsForumReadable(r.Context(), forumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}
		if !readable {
			ep_util.RenderErrorHttp404NotFoundWithMessage("Page not found.", w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["MessageTitle"] = "Forum"
		templateCtx["MessageBody"] = "This page is under construction."
		ep_util.RenderTemplate("view_simple_message", templateCtx, w, r)
	})
	mux.Handle("GET /forum/index/{forumId}", getForumIndexHandler)
}
