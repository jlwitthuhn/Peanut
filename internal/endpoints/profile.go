// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package endpoints

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/service"
)

func RegisterProfileHandlers(mux *http.ServeMux, userService service.UserService) {
	getProfileHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		user, err := userService.GetUserRowById(r, id)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to query user details from database", w, r)
			return
		}

		if user == nil {
			ep_util.RenderErrorHttp404NotFoundWithMessage("Failed to find specified user", w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["ProfileId"] = user.Id
		ep_util.RenderTemplate("_profile", templateCtx, w, r)
	})
	mux.Handle("GET /profile/{id}", getProfileHandler)
}
