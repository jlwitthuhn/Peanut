// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package endpoints

import (
	"html/template"
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/keynames/configkey"
	"peanut/internal/logger"
	"peanut/internal/msgfmt"
	"peanut/internal/service"
)

func RegisterIndexHandlers(mux *http.ServeMux, configService service.ConfigService) {

	getIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		welcomeMessage, err := configService.GetString(r.Context(), configkey.StringWelcomeMessage)
		if err != nil {
			logger.Error(r.Context(), "Error retrieving welcome message, using error message.", err)
		}
		var welcomeMessageText string
		if err != nil || welcomeMessage == nil {
			welcomeMessageText = "Error: unable to retrieve welcome message."
		} else {
			welcomeMessageText = *welcomeMessage
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["WelcomeMessage"] = template.HTML(msgfmt.Format(welcomeMessageText))
		ep_util.RenderTemplate("view_index", templateCtx, w, r)
	})
	mux.Handle("GET /{$}", getIndexHandler)
}
