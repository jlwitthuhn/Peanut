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
	"peanut/internal/template"
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

		theTemplate := template.GetTemplate("_admin/users")
		err = theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing template:", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to execute template.", w, r)
			return
		}
	})
	mux.Handle("GET /admin/users", getHandler)
}
