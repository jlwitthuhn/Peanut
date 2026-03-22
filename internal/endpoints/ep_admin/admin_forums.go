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
	"strconv"
)

func registerAdminForumsHandlers(mux *http.ServeMux, forumsService service.ForumsService) {
	getSectionsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		ep_util.RenderTemplate("_admin/forum/sections", templateCtx, w, r)
	})
	mux.Handle("GET /admin/forum/sections", getSectionsHandler)

	getSectionsAddHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		ep_util.RenderTemplate("_admin/forum/sections/add", templateCtx, w, r)
	})
	mux.Handle("GET /admin/forum/sections/add", getSectionsAddHandler)

	postSectionsAddHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		title := r.PostFormValue("title")
		orderStr := r.PostFormValue("order")

		ordering, err := strconv.Atoi(orderStr)
		if err != nil {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Order must be a valid integer.", w, r)
			return
		}

		err = forumsService.CreateSection(r, title, ordering)
		if err != nil {
			logger.Error(r, "Failed to create forum section: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create forum section.", w, r)
			return
		}

		err = ep_util.CommitTransactionForRequest(r)
		if err != nil {
			logger.Error(r, "Failed to commit transaction: ", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		RenderSimpleAdminMessage("Success", "Forum section '"+title+"' has been created.", w, r)
	})
	mux.Handle("POST /admin/forum/sections/add", postSectionsAddHandler)
}
