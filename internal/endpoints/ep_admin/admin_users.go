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

func registerAdminUsersHandlers(mux *http.ServeMux, groupService service.GroupService) {
	getHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		groupList, err := groupService.GetAllGroupNames(r)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to query group names.", w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["Groups"] = groupList
		ep_util.RenderTemplate("_admin/users", templateCtx, w, r)
	})
	mux.Handle("GET /admin/users", getHandler)
}
