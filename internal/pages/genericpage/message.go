// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package genericpage

import (
	"net/http"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/template"
)

func RenderSimpleMessage(title string, message string, w http.ResponseWriter, r *http.Request) {
	templateCtx := make(map[string]any)
	templateCtx["MessageBody"] = message
	templateCtx["MessageTitle"] = title
	templateCtx["RequestDuration"] = middleutil.RequestTimerFinish(r)

	theTemplate := template.GetTemplate("_simple_message")
	err := theTemplate.Execute(w, templateCtx)
	if err != nil {
		logger.Error("Error executing simple message template:", err)
	}
}
