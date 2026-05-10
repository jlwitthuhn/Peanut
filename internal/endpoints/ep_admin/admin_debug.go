// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"fmt"
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/service"
	"time"
)

func registerAdminDebugHandlers(mux *http.ServeMux, forumService service.ForumService, forumThreadhreadService service.ForumThreadService) {
	getDataHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		ep_util.RenderTemplate("view_admin/debug", templateCtx, w, r)
	})
	mux.Handle("GET /admin/debug/data", getDataHandler)

	postDataCreateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		section1Id, err := forumService.CreateSection(ctx, "Section 1", 1)
		if err != nil {
			logger.Error(ctx, "Failed to create Section 1: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Section 1.", w, r)
			return
		}

		forum1Id, err := forumService.CreateForum(ctx, section1Id, "Forum 1", "The first forum", 1, "Public")
		if err != nil {
			logger.Error(ctx, "Failed to create Forum 1: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Forum 1.", w, r)
			return
		}

		_, err = forumService.CreateForum(ctx, section1Id, "Forum 2", "The second forum", 2, "Public")
		if err != nil {
			logger.Error(ctx, "Failed to create Forum 2: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Forum 2.", w, r)
			return
		}

		section2Id, err := forumService.CreateSection(ctx, "Section 2", 2)
		if err != nil {
			logger.Error(ctx, "Failed to create Section 2: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Section 2.", w, r)
			return
		}

		_, err = forumService.CreateForum(ctx, section2Id, "Forum 3", "This is also a forum", 1, "Public")
		if err != nil {
			logger.Error(ctx, "Failed to create Forum 3: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Forum 3.", w, r)
			return
		}

		message1 := `
# Title
This is a thread for testing threads.

## Subtitle
*Markdown* formatting works here. This thread has **many** replies.
`
		thread1Id, err := forumThreadhreadService.CreateThread(ctx, forum1Id, "Test Thread 1", message1)
		if err != nil {
			logger.Error(ctx, "Failed to create test thread 1: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create test thread 1.", w, r)
			return
		}

		thread2Id, err := forumThreadhreadService.CreateThread(ctx, forum1Id, "Test Thread 2", "This is another test thread.")
		if err != nil {
			logger.Error(ctx, "Failed to create test thread 2: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("This thread only has a few replies.", w, r)
			return
		}

		for i := 1; i <= 100; i++ {
			time.Sleep(time.Millisecond)
			_, err = forumThreadhreadService.AddThreadPost(ctx, thread1Id, fmt.Sprintf("Reply %d", i))
			if err != nil {
				logger.Error(ctx, "Failed to add reply to thread1: ", err)
				ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to add reply thread1.", w, r)
				return
			}
			if i < 10 {
				_, err = forumThreadhreadService.AddThreadPost(ctx, thread2Id, fmt.Sprintf("Reply %d", i))
				if err != nil {
					logger.Error(ctx, "Failed to add reply to thread2: ", err)
					ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to add reply thread2.", w, r)
					return
				}
			}
		}

		err = ep_util.CommitTransactionForRequest(r)
		if err != nil {
			logger.Error(ctx, "Failed to commit transaction: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		RenderSimpleAdminMessage("Success", "Debug data has been created.", w, r)
	})
	mux.Handle("POST /admin/debug/data/create", postDataCreateHandler)
}
