// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package pages

import (
	"net/http"
	"peanut/internal/logger"
	"peanut/internal/template"
)

func RegisterIndexHandlers(mux *http.ServeMux) {

	indexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		theTemplate := template.GetTemplate("_index")
		err := theTemplate.Execute(w, nil)
		if err != nil {
			logger.Error("Error executing template:", err)
		}
	})
	mux.Handle("/{$}", indexHandler)
}
