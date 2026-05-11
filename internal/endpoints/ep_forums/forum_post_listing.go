// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"fmt"
	"html/template"
	"net/http"
	"peanut/internal/data"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/msgfmt"
	"peanut/internal/service"
	"strconv"
)

type forumPostListingItem struct {
	data.ForumPostListingViewRow
	FormattedMessage template.HTML
}

func registerForumPostListingHandlers(mux *http.ServeMux, forumsService service.ForumService, forumThreadService service.ForumThreadService) {
	pageUrl := func(threadId string, p int) string {
		if p == 1 {
			return "/forum/thread/" + threadId
		}
		return fmt.Sprintf("/forum/thread/%s/page/%d", threadId, p)
	}

	renderPage := func(w http.ResponseWriter, r *http.Request, threadId string, pageNum int) {
		thread, err := forumThreadService.GetThreadRowById(r.Context(), threadId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load thread.", w, r)
			return
		}
		if thread == nil {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}
		if thread.Visibility != "Public" {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		forum, err := forumsService.GetForumRowById(r.Context(), thread.ForumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load forum.", w, r)
			return
		}

		readable, err := forumsService.CanReadForum(r.Context(), thread.ForumId)
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

		page, err := forumThreadService.GetForumPostListing(r.Context(), threadId, pageNum)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to load posts.", w, r)
			return
		}
		if page == nil {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}

		postItems := make([]forumPostListingItem, len(page.Posts))
		for i, post := range page.Posts {
			postItems[i] = forumPostListingItem{
				ForumPostListingViewRow: post,
				FormattedMessage:        template.HTML(msgfmt.Format(post.PostMessage)),
			}
		}

		canReply, err := forumThreadService.CanPostReplyInThread(r.Context(), thread.ForumId)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to check permissions.", w, r)
			return
		}

		pagination := ep_util.BuildPagination(page.CurrentPage, page.TotalPages, func(p int) string {
			return pageUrl(thread.Id, p)
		})

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["CanReply"] = canReply
		templateCtx["ThreadId"] = thread.Id
		templateCtx["ThreadTitle"] = thread.Title
		templateCtx["ForumId"] = forum.Id
		templateCtx["ForumName"] = forum.Name
		templateCtx["SectionId"] = section.Id
		templateCtx["SectionName"] = section.Name
		templateCtx["Posts"] = postItems
		templateCtx["Pagination"] = pagination
		ep_util.RenderTemplate("view_forum/post_listing", templateCtx, w, r)
	}

	getPostListingHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		threadId := r.PathValue("threadId")
		renderPage(w, r, threadId, 1)
	})
	mux.Handle("GET /forum/thread/{threadId}", getPostListingHandler)

	getPostListingPageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		threadId := r.PathValue("threadId")
		pageStr := r.PathValue("pageNum")
		pageNum, err := strconv.Atoi(pageStr)
		if err != nil || pageNum < 1 {
			ep_util.RenderErrorHttp404NotFound(w, r)
			return
		}
		if pageNum == 1 {
			http.Redirect(w, r, pageUrl(threadId, 1), http.StatusMovedPermanently)
			return
		}
		renderPage(w, r, threadId, pageNum)
	})
	mux.Handle("GET /forum/thread/{threadId}/page/{pageNum}", getPostListingPageHandler)
}
