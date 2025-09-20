// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package pages

import (
	"net/http"
	"peanut/internal/database"
	"peanut/internal/database/data/data_config"
	"peanut/internal/database/data/data_meta"
	"peanut/internal/logger"
	"peanut/internal/middleutil"
	"peanut/internal/pages/genericpage"
	"peanut/internal/service/db_service"
	"peanut/internal/template"
)

func isEmailValid(email string) bool {
	return len(email) > 1
}

func isPasswordValid(password string) bool {
	return len(password) > 1
}

func isUsernameValid(username string) bool {
	return len(username) > 1
}

func RegisterSetupHandlers(mux *http.ServeMux) {

	getSetupHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := make(map[string]any)
		templateCtx["RequestDuration"] = middleutil.RequestTimerFinish(r)

		theTemplate := template.GetTemplate("_setup")
		err := theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error("Error executing setup template:", err)
		}
	})
	mux.Handle("GET /setup", getSetupHandler)

	postSetupHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		configTableExists, configTableErr := db_service.DoesTableExist("config_int")
		if configTableErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to query database.", w, r)
			return
		}
		if configTableExists {
			w.WriteHeader(http.StatusConflict)
			genericpage.RenderSimpleMessage("Error", "Database has already been initialized.", w, r)
			return
		}

		adminName := r.PostFormValue("admin-name")
		email := r.PostFormValue("email")
		adminPassword := r.PostFormValue("admin-pass")
		adminPassword2 := r.PostFormValue("admin-pass-2")

		if !isUsernameValid(adminName) {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "Username is not valid.", w, r)
			return
		}
		if !isEmailValid(email) {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "Email is not valid.", w, r)
			return
		}
		if adminPassword != adminPassword2 {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "Passwords must match.", w, r)
			return
		}
		if !isPasswordValid(adminPassword) {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "Passwords is not valid.", w, r)
			return
		}

		tx, txErr := database.PostgresHandle().BeginTx(r.Context(), nil)
		if txErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to start database transaction.", w, r)
			return
		}
		defer tx.Rollback()

		metaErr := data_meta.CreateDBObjects(tx)
		if metaErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to create meta db objects.", w, r)
			return
		}
		configErr := data_config.CreateDBObjects(tx)
		if configErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to create config db objects.", w, r)
			return
		}

		commitErr := tx.Commit()
		if commitErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to commit database transaction.", w, r)
			return
		}

		genericpage.RenderSimpleMessage("Complete", "Peanut has been initialized.", w, r)
	})
	mux.Handle("POST /setup", postSetupHandler)
}
