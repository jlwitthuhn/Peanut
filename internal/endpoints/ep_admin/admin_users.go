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

func registerAdminUsersHandlers(mux *http.ServeMux, groupService service.GroupService, userService service.UserService) {
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

	getListAllHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userRows, err := userService.GetUserRowsAll(r)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to query user names.", w, r)
			return
		}

		// Truncate IDs so they don't take up the whole width of the page
		for i := range userRows {
			userRows[i].Id = userRows[i].Id[:13] + "..."
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["Users"] = userRows
		ep_util.RenderTemplate("_admin/users_list", templateCtx, w, r)
	})
	mux.Handle("GET /admin/users/list/all", getListAllHandler)

	postListByNamePatternHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pattern := r.PostFormValue("pattern")
		userRows, err := userService.GetUserRowsLikeName(r, pattern)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to query user names.", w, r)
		}

		// Truncate IDs so they don't take up the whole width of the page
		for i := range userRows {
			userRows[i].Id = userRows[i].Id[:13] + "..."
		}
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["Users"] = userRows
		ep_util.RenderTemplate("_admin/users_list", templateCtx, w, r)
	})
	mux.Handle("POST /admin/users/list/by_name", postListByNamePatternHandler)
}
