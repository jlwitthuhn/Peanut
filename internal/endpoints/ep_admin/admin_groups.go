// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import "net/http"

func registerAdminGroupsHandlers(mux *http.ServeMux) {
	getHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RenderSimpleAdminMessage("Groups", "View and edit group information here.", w, r)
	})
	mux.Handle("GET /admin/groups", getHandler)
}
