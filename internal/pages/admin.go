// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package pages

import (
	"net/http"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/pages/genericpage"
	"peanut/internal/pages/templatecontext"
	"peanut/internal/perms"
	"peanut/internal/template"
)

func RegisterAdminHandlers(mux *http.ServeMux) {
	getIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleutil.RequestHasPermission(r, perms.Admin_Gui_View) == false {
			genericpage.RenderErrorHttp403Forbidden(w, r)
			return
		}
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		theTemplate := template.GetTemplate("_admin/index")
		err := theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing template:", err)
		}
	})
	mux.Handle("GET /admin", getIndexHandler)
}
