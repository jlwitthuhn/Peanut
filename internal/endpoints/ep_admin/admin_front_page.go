// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_admin

import (
	"net/http"
	"peanut/internal/data/configkey"
	"peanut/internal/endpoints/genericpage"
	"peanut/internal/endpoints/requtil"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/security/perms"
	"peanut/internal/service"
	"peanut/internal/template"
)

func registerAdminFrontPageHandlers(mux *http.ServeMux, configService service.ConfigService) {
	getFrontPageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleutil.RequestHasPermission(r, perms.Admin_FrontPage_Edit) == false {
			genericpage.RenderErrorHttp403Forbidden(w, r)
			return
		}

		welcomeMessage, err := configService.GetString(r, configkey.StringWelcomeMessage)
		if err != nil {
			logger.Error(r, "Error retrieving welcome message.", err)
			genericpage.RenderErrorHttp403Forbidden(w, r)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		templateCtx["WelcomeMessage"] = welcomeMessage
		theTemplate := template.GetTemplate("_admin/front_page")
		err = theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing template:", err)
			genericpage.RenderErrorHttp500InternalServerErrorWithMessage("Failed to execute template.", w, r)
			return
		}
	})
	mux.Handle("GET /admin/front_page", getFrontPageHandler)

	postFrontPageWelcomeMessageHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if middleutil.RequestHasPermission(r, perms.Admin_FrontPage_Edit) == false {
			genericpage.RenderErrorHttp403Forbidden(w, r)
			return
		}

		newMessage := r.PostFormValue("message")
		confirm := r.PostFormValue("confirm")

		if confirm != "on" {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "You must check the 'Confirm' box to set the welcome message.", w, r)
			return
		}

		err := configService.SetString(r, configkey.StringWelcomeMessage, newMessage)
		if err != nil {
			logger.Error(r, "Failed to set new welcome message.", err)
			genericpage.RenderErrorHttp500InternalServerErrorWithMessage("Failed to set new welcome message.", w, r)
			return
		}

		err = requtil.CommitTransactionForRequest(r)
		if err != nil {
			genericpage.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		genericpage.RenderSimpleMessage("Success", "New welcome message has been set: "+newMessage, w, r)
	})
	mux.Handle("POST /admin/front_page/welcome_message", postFrontPageWelcomeMessageHandler)
}
