// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package pages

import (
	"net/http"
	"peanut/internal/data/configkey"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/pages/genericpage"
	"peanut/internal/pages/templatecontext"
	"peanut/internal/security/perms"
	"peanut/internal/service"
	"peanut/internal/template"
	"runtime"
	"time"
)

func RegisterAdminHandlers(mux *http.ServeMux, configService service.ConfigService, databaseService service.DatabaseService) {
	getIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleutil.RequestHasPermission(r, perms.Admin_Gui_View) == false {
			genericpage.RenderErrorHttp403Forbidden(w, r)
			return
		}

		initTime, initTimeErr := configService.GetInt(nil, configkey.IntInitializedTime)
		if initTimeErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to query init time.", w, r)
			return
		}

		dbVersion, dbVersionErr := databaseService.GetPostgresVersion()
		if dbVersionErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to query postgres version.", w, r)
			return
		}

		var websiteInfo = make(map[string]string)
		websiteInfo["Initialized Time"] = time.Unix(initTime, 0).UTC().Format("2006-01-02 15:04:05 MST")

		var envInfo = make(map[string]string)
		envInfo["Go Runtime"] = runtime.Version()
		envInfo["PostgreSQL Version"] = dbVersion
		envInfo["Host OS"] = runtime.GOOS
		envInfo["Host Arch"] = runtime.GOARCH

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["WebsiteInfo"] = websiteInfo
		templateCtx["EnvironmentInfo"] = envInfo

		theTemplate := template.GetTemplate("_admin/index")
		err := theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing template:", err)
		}
	})
	mux.Handle("GET /admin", getIndexHandler)
}
