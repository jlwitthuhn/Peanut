// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package genericpage

import (
	"net/http"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/template"
)

func RenderErrorHttp403Forbidden(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusForbidden)
	RenderSimpleMessage("403 - Forbidden", "You do not have permission to access this page.", w, r)
}

func RenderErrorHttp500InternalServerError(w http.ResponseWriter, r *http.Request) {
	RenderErrorHttp500InternalServerErrorWithMessage("The server encountered an error while processing your request.", w, r)
}

func RenderErrorHttp500InternalServerErrorWithMessage(message string, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	RenderSimpleMessage("500 - Internal Server Error", message, w, r)
}

func RenderSimpleMessage(title string, message string, w http.ResponseWriter, r *http.Request) {
	templateCtx := templatecontext.GetStandardTemplateContext(r)
	templateCtx["MessageBody"] = message
	templateCtx["MessageTitle"] = title

	theTemplate := template.GetTemplate("_simple_message")
	err := theTemplate.Execute(w, templateCtx)
	if err != nil {
		logger.Error(r, "Error executing simple message template:", err)
	}
}
