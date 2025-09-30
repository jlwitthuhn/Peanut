// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package pages

import (
	"net/http"
	"peanut/internal/data/configkey"
	"peanut/internal/logger"
	"peanut/internal/pages/templatecontext"
	"peanut/internal/service"
	"peanut/internal/template"
)

func RegisterIndexHandlers(mux *http.ServeMux, configService service.ConfigService) {

	getIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		welcomeMessage, err := configService.GetString(nil, configkey.StringWelcomeMessage)
		if err != nil {
			logger.Error(r, "Error retrieving welcome message, using error message.", err)
			welcomeMessage = "Error: unable to retrieve welcome message."
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["WelcomeMessage"] = welcomeMessage
		theTemplate := template.GetTemplate("_index")
		err = theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing template:", err)
		}
	})
	mux.Handle("GET /{$}", getIndexHandler)
}
