// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import "net/http"

func registerAdminScheduledJobHandlers(mux *http.ServeMux) {
	getScheduledJobHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		RenderSimpleAdminMessage("TODO", "Page incomplete.", w, r)
	})
	mux.HandleFunc("/admin/scheduled_jobs", getScheduledJobHandler)
}
