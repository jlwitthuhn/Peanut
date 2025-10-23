// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/security/perms"
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
	mux.HandleFunc("GET /admin/scheduled_jobs", getScheduledJobHandler)

	postScheduledJobRunHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleutil.RequestHasPermission(r, perms.Admin_ScheduledJob_Run) == false {
			ep_util.RenderErrorHttp403Forbidden(w, r)
			return
		}

		id := r.PostFormValue("id")
		_, err := scheduledJobService.GetJobNameById(r, id)
		if err != nil {
			logger.Error(r, "Failed to query scheduled job:", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to query scheduled job.", w, r)
			return
		}

		RenderSimpleAdminMessage("TODO", "Not implemented.", w, r)
	})
	mux.HandleFunc("POST /admin/scheduled_jobs/run", postScheduledJobRunHandler)
}
