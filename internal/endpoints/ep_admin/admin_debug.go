// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/service"
)

func registerAdminDebugHandlers(mux *http.ServeMux, forumsService service.ForumService) {
	getDataHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		ep_util.RenderTemplate("view_admin/debug", templateCtx, w, r)
	})
	mux.Handle("GET /admin/debug/data", getDataHandler)

	postDataCreateHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		section1Id, err := forumsService.CreateSection(ctx, "Section 1", 1)
		if err != nil {
			logger.Error(ctx, "Failed to create Section 1: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Section 1.", w, r)
			return
		}

		_, err = forumsService.CreateForum(ctx, section1Id, "Forum 1", "", 1, "Public")
		if err != nil {
			logger.Error(ctx, "Failed to create Forum 1: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Forum 1.", w, r)
			return
		}

		_, err = forumsService.CreateForum(ctx, section1Id, "Forum 2", "", 2, "Public")
		if err != nil {
			logger.Error(ctx, "Failed to create Forum 2: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Forum 2.", w, r)
			return
		}

		section2Id, err := forumsService.CreateSection(ctx, "Section 2", 2)
		if err != nil {
			logger.Error(ctx, "Failed to create Section 2: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Section 2.", w, r)
			return
		}

		_, err = forumsService.CreateForum(ctx, section2Id, "Forum 3", "", 1, "Public")
		if err != nil {
			logger.Error(ctx, "Failed to create Forum 3: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create Forum 3.", w, r)
			return
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
