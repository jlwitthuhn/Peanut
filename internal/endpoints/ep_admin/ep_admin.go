// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/middleware"
	"peanut/internal/security/perms"
	"peanut/internal/service"
)

func RegisterAdminHandlers(mux *http.ServeMux, configService service.ConfigService, databaseService service.DatabaseService, sessionService service.SessionService, userService service.UserService) {
	adminMux := http.NewServeMux()
	registerAdminIndexHandlers(adminMux, configService, databaseService, sessionService, userService)
	registerAdminFrontPageHandlers(adminMux, configService)

	wrappedAdminMux := middleware.WrapHandler(adminMux, middleware.CheckPermissions(perms.Admin_Gui_View))

	mux.Handle("/admin", wrappedAdminMux)
	mux.Handle("/admin/", wrappedAdminMux)
}
