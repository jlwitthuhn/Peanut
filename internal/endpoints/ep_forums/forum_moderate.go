// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/security/perms"
	"peanut/internal/service"
)

func registerForumModerateHandlers(mux *http.ServeMux, forumThreadService service.ForumThreadService) {
	getModerateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ep_util.RequirePermissionOr403(w, r, perms.Forum_Moderate) {
			return
		}

		threadId := r.PathValue("threadId")

		thread, err := forumThreadService.GetThreadRowById(r.Context(), threadId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load thread.", w, r)
			return
		}
		if thread == nil {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["ThreadId"] = thread.Id
		templateCtx["Title"] = thread.Title
		templateCtx["Pinned"] = thread.Pinned
		templateCtx["Locked"] = thread.Locked
		ep_util.RenderTemplate("view_forum/moderate", templateCtx, w, r)
	})
	mux.Handle("GET /forum/thread/{threadId}/moderate", getModerateHandler)

	postModerateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !ep_util.RequirePermissionOr403(w, r, perms.Forum_Moderate) {
			return
		}

		threadId := r.PathValue("threadId")

		thread, err := forumThreadService.GetThreadRowById(r.Context(), threadId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load thread.", w, r)
			return
		}
		if thread == nil {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		title := r.PostFormValue("title")
		if title == "" {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Thread title is required.", w, r)
			return
		}
		pinned := r.PostFormValue("pinned") != ""
		locked := r.PostFormValue("locked") != ""

		err = forumThreadService.ModerateThread(r.Context(), thread.Id, title, pinned, locked)
		if err != nil {
			logger.Error(r.Context(), "Failed to moderate thread: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to update thread.", w, r)
			return
		}

		err = ep_util.CommitTransactionForRequest(r)
		if err != nil {
			logger.Error(r.Context(), "Failed to commit transaction: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		http.Redirect(w, r, "/forum/thread/"+thread.Id, http.StatusSeeOther)
	})
	mux.Handle("POST /forum/thread/{threadId}/moderate", postModerateHandler)
}
