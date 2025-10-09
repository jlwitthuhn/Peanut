// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/logger"
	"peanut/internal/service"
)

func DatabaseInitCheck(dbService service.DatabaseService, setupHandler http.Handler) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tableExists, err := dbService.DoesTableExist(r, "config_int")
			if err != nil {
				logger.Fatal(r, err)
			}
			if !tableExists {
				// Allow access to setup page only when DB is not initialized
				if r.URL.Path == "/setup" {
					setupHandler.ServeHTTP(w, r)
				} else {
					w.WriteHeader(http.StatusNotFound)
					ep_util.RenderSimpleMessage("Database Not Initialized", "The data must be configured before Peanut can be used.", w, r)
				}
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
