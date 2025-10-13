// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/middleware"
	"peanut/internal/security/perms"
	"peanut/internal/service"
)

func RegisterAdminHandlers(
	mux *http.ServeMux,
	configService service.ConfigService,
	databaseService service.DatabaseService,
	groupService service.GroupService,
	sessionService service.SessionService,
	userService service.UserService,
) {
	adminMux := http.NewServeMux()
	registerAdminIndexHandlers(adminMux, configService, databaseService, sessionService, userService)
	registerAdminFrontPageHandlers(adminMux, configService)
	registerAdminGroupsHandlers(adminMux, groupService)
	registerAdminScheduledJobHandlers(adminMux)
	registerAdminUsersHandlers(adminMux, groupService, userService)

	wrappedAdminMux := middleware.WrapHandler(adminMux, middleware.CheckPermissions(perms.Admin_Gui_View))

	mux.Handle("/admin", wrappedAdminMux)
	mux.Handle("/admin/", wrappedAdminMux)
}

func RenderSimpleAdminMessage(title string, message string, w http.ResponseWriter, r *http.Request) {
	templateCtx := templatecontext.GetStandardTemplateContext(r)
	templateCtx["MessageBody"] = message
	templateCtx["MessageTitle"] = title
	ep_util.RenderTemplate("_admin/simple_message", templateCtx, w, r)
}
