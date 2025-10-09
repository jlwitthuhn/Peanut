// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package endpoints

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/service"
	"peanut/internal/template"
	"peanut/internal/validator"
)

func RegisterSetupHandlers(mux *http.ServeMux, dbService service.DatabaseService, setupService service.SetupService) {

	getSetupHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		theTemplate := template.GetTemplate("_setup")
		err := theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing setup template:", err)
		}
	})
	mux.Handle("GET /setup", getSetupHandler)

	postSetupHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		configTableExists, configTableErr := dbService.DoesTableExist(r, "config_int")
		if configTableErr != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to query data.", w, r)
			return
		}
		if configTableExists {
			w.WriteHeader(http.StatusConflict)
			ep_util.RenderSimpleMessage("Error", "Database has already been initialized.", w, r)
			return
		}

		adminName := r.PostFormValue("admin-name")
		email := r.PostFormValue("email")
		adminPassword := r.PostFormValue("admin-pass")
		adminPassword2 := r.PostFormValue("admin-pass-2")

		if err := validator.ValidateUsername(adminName); err != nil {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Username is not valid: "+err.Error(), w, r)
			return
		}
		if err := validator.ValidateEmail(email); err != nil {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Email is not valid: "+err.Error(), w, r)
			return
		}
		if adminPassword != adminPassword2 {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Passwords must match.", w, r)
			return
		}
		if err := validator.ValidatePassword(adminPassword); err != nil {
			ep_util.RenderErrorHttp400BadRequestWithMessage("Password is not valid: "+err.Error(), w, r)
			return
		}

		logger.Info(r, "Input valid, initializing...")
		err := setupService.InitializeDatabase(r, adminName, email, adminPassword)
		if err != nil {
			logger.Error(r, "Error initializing database:", err)
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to initialize database.", w, r)
			return
		}

		logger.Info(r, "Committing transaction...")

		err = ep_util.CommitTransactionForRequest(r)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		logger.Info(r, "Peanut initialization complete.")
		ep_util.RenderSimpleMessage("Complete", "Peanut has been initialized.", w, r)
	})
	mux.Handle("POST /setup", postSetupHandler)
}
