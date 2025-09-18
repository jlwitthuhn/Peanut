// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package pages

import (
	"net/http"
	"time"

	"peanut/internal/logger"
	"peanut/internal/middleware"
	"peanut/internal/template"
)

func RegisterIndexHandlers(mux *http.ServeMux) {

	getIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestBegin := r.Context().Value(middleware.RequestTimerBeginKey).(time.Time)
		requestDurationUs := time.Now().Sub(requestBegin).Microseconds()
		requestDurationUs -= requestDurationUs % 10 // Two decimal places in milliseconds
		requestDurationMs := float64(requestDurationUs) / 1000.0

		templateCtx := make(map[string]any)
		templateCtx["RequestDuration"] = requestDurationMs
		theTemplate := template.GetTemplate("_index")

		err := theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error("Error executing template:", err)
		}
	})
	mux.Handle("GET /{$}", getIndexHandler)
}
