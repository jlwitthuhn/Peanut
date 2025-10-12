// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/service"
)

func registerAdminGroupsHandlers(mux *http.ServeMux, groupService service.GroupService) {
	getHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		ep_util.RenderTemplate("_admin/groups", templateCtx, w, r)
	})
	mux.Handle("GET /admin/groups", getHandler)

	getListAllHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupRows, err := groupService.GetAllGroupRows(r)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to query group list.", w, r)
			return
		}

		// Truncate IDs so they don't take up the whole width of the page
		for i := range groupRows {
			groupRows[i].Id = groupRows[i].Id[:13] + "..."
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["Groups"] = groupRows
		ep_util.RenderTemplate("_admin/groups_list", templateCtx, w, r)
	})
	mux.Handle("GET /admin/groups/list/all", getListAllHandler)
}
