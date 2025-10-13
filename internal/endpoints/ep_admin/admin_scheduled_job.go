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
)

func registerAdminScheduledJobHandlers(mux *http.ServeMux, scheduledJobService service.ScheduledJobService) {
	getScheduledJobHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		jobs, err := scheduledJobService.GetAllJobSummaries(r)
		if err != nil {
			logger.Error(r, "Failed to query scheduled job summaries:", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to query scheduled job summaries.", w, r)
			return
		}
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["Jobs"] = jobs
		ep_util.RenderTemplate("_admin/scheduled_jobs", templateCtx, w, r)
	})
	mux.HandleFunc("/admin/scheduled_jobs", getScheduledJobHandler)
}
