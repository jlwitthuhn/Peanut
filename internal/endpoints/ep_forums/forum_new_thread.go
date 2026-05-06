// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/service"
)

func registerForumNewThreadHandlers(mux *http.ServeMux, forumsService service.ForumService, forumThreadService service.ForumThreadService) {
	getNewThreadHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		forumId := r.PathValue("forumId")

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

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["ForumId"] = forum.Id
		templateCtx["ForumName"] = forum.Name
		templateCtx["SectionId"] = section.Id
		templateCtx["SectionName"] = section.Name
		ep_util.RenderTemplate("view_forum/new_thread", templateCtx, w, r)
	})
	mux.Handle("GET /forum/new_thread/{forumId}", getNewThreadHandler)

	postNewThreadHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		urlForumId := r.PathValue("forumId")
		formForumId := r.PostFormValue("forum_id")

		if urlForumId != formForumId {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Forum ID in URL does not match forum ID in form.", w, r)
			return
		}

		forum, err := forumsService.GetForumRowById(r.Context(), urlForumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}
		if forum == nil {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		writable, err := forumThreadService.CanPostThreadInForum(r.Context(), urlForumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}
		if !writable {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		title := r.PostFormValue("title")
		if title == "" {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Thread title is required.", w, r)
			return
		}

		message := r.PostFormValue("message")
		if message == "" {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Thread message is required.", w, r)
			return
		}

		_, err = forumThreadService.CreateThread(r.Context(), urlForumId, title, message)
		if err != nil {
			logger.Error(r.Context(), "Failed to create thread: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create thread.", w, r)
			return
		}

		err = ep_util.CommitTransactionForRequest(r)
		if err != nil {
			logger.Error(r.Context(), "Failed to commit transaction: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		http.Redirect(w, r, "/forum/index/"+urlForumId, http.StatusSeeOther)
	})
	mux.Handle("POST /forum/new_thread/{forumId}", postNewThreadHandler)
}
