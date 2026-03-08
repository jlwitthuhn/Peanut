// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
)

func registerAdminForumsHandlers(mux *http.ServeMux) {
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
		ep_util.RenderSimpleMessage("Not Implemented", "This endpoint has not been implemented yet.", w, r)
	})
	mux.Handle("POST /admin/forum/sections/add", postSectionsAddHandler)
}
