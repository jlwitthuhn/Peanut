// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_util

import (
	"net/http"
	"peanut/internal/logger"
	"peanut/internal/template"
)

func RenderTemplate(templateName string, context map[string]any, w http.ResponseWriter, r *http.Request) {
	theTemplate := template.GetTemplate(templateName)
	err := theTemplate.Execute(w, context)
	if err != nil {
		logger.Error(r, "Error executing template:", err)
		RenderErrorHttp500InternalServerErrorWithMessage("Failed to execute template.", w, r)
		return
	}
}
