// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package endpoints

import (
	"net/http"
	"peanut/internal/data/configkey"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/service"
)

func RegisterIndexHandlers(mux *http.ServeMux, configService service.ConfigService) {

	getIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		welcomeMessage, err := configService.GetString(r, configkey.StringWelcomeMessage)
		if err != nil {
			logger.Error(r, "Error retrieving welcome message, using error message.", err)
			welcomeMessage = "Error: unable to retrieve welcome message."
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["WelcomeMessage"] = welcomeMessage
		ep_util.RenderTemplate("_index", templateCtx, w, r)
	})
	mux.Handle("GET /{$}", getIndexHandler)
}
