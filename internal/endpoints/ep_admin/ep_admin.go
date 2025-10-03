// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/service"
)

func RegisterAdminHandlers(mux *http.ServeMux, configService service.ConfigService, databaseService service.DatabaseService, userService service.UserService) {
	registerAdminIndexHandlers(mux, configService, databaseService, userService)
	registerAdminFrontPageHandlers(mux, configService)
}
