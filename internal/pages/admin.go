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
	"strconv"
	"time"
)

type adminPageStringPair struct {
	A string
	B string
}

func RegisterAdminHandlers(mux *http.ServeMux, configService service.ConfigService, databaseService service.DatabaseService, userService service.UserService) {
	getIndexHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleutil.RequestHasPermission(r, perms.Admin_Gui_View) == false {
			genericpage.RenderErrorHttp403Forbidden(w, r)
			return
		}

		initTime, err := configService.GetInt(nil, configkey.IntInitializedTime)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to query init time.", w, r)
			return
		}

		dbVersion, err := databaseService.GetPostgresVersion()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to query postgres version.", w, r)
			return
		}

		userCount, err := userService.CountUsers(nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to query user count.", w, r)
			return
		}

		userSessionCount, err := userService.CountUsersWithValidSession(nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to query user session count.", w, r)
			return
		}

		var websiteInfo = []adminPageStringPair{
			{A: "Initialized time", B: time.Unix(initTime, 0).UTC().Format("2006-01-02 15:04:05 MST")},
			{A: "Registered users", B: strconv.FormatInt(userCount, 10)},
			{A: "Logged in users", B: strconv.FormatInt(userSessionCount, 10)},
		}

		var envInfo = []adminPageStringPair{
			{A: "Go runtime", B: runtime.Version()},
			{A: "PostgreSQL version", B: dbVersion},
			{A: "Host OS", B: runtime.GOOS},
			{A: "Host arch", B: runtime.GOARCH},
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["WebsiteInfo"] = websiteInfo
		templateCtx["EnvironmentInfo"] = envInfo

		theTemplate := template.GetTemplate("_admin/index")
		err = theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing template:", err)
		}
	})
	mux.Handle("GET /admin", getIndexHandler)
}
