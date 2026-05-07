// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"net/http"
	"peanut/internal/data"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/service"
)

func registerForumReplyHandlers(mux *http.ServeMux, forumsService service.ForumService, forumThreadService service.ForumThreadService) {
	enforceReplyAccess := func(w http.ResponseWriter, r *http.Request) *data.ForumThreadRow {
		threadId := r.PathValue("threadId")

		thread, err := forumThreadService.GetThreadRowById(r.Context(), threadId)
		if err != nil || thread == nil || thread.Visibility != "Public" {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return nil
		}

		canReply, err := forumThreadService.CanPostReplyInThread(r.Context(), thread.ForumId)
		if err != nil || !canReply {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return nil
		}

		return thread
	}

	getReplyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thread := enforceReplyAccess(w, r)
		if thread == nil {
			return
		}

		forum, err := forumsService.GetForumRowById(r.Context(), thread.ForumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
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

	postReplyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		thread := enforceReplyAccess(w, r)
		if thread == nil {
			return
		}

		message := r.PostFormValue("message")
		if message == "" {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Message is required.", w, r)
			return
		}

		_, err := forumThreadService.AddThreadPost(r.Context(), thread.Id, message)
		if err != nil {
			logger.Error(r.Context(), "Failed to add reply: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to add reply.", w, r)
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
	mux.Handle("POST /forum/thread/{threadId}/reply", postReplyHandler)
}
