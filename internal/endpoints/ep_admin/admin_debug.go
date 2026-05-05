// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
)

func registerAdminDebugHandlers(mux *http.ServeMux) {
	getDataHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		ep_util.RenderTemplate("view_admin/debug", templateCtx, w, r)
	})
	mux.Handle("GET /admin/debug/data", getDataHandler)
}
