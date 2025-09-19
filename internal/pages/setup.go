// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package pages

import (
	"net/http"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/template"
)

func RegisterSetupHandlers(mux *http.ServeMux) {

	getSetupHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := make(map[string]any)
		templateCtx["RequestDuration"] = middleutil.RequestTimerFinish(r)

		theTemplate := template.GetTemplate("_setup")
		err := theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error("Error executing setup template:", err)
		}
	})
	mux.Handle("GET /setup", getSetupHandler)
}
