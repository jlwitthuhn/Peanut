// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package endpoints

import (
	"net/http"
	"peanut/internal/endpoints/genericpage"
	"peanut/internal/endpoints/requtil"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/security/perms/permgroups"
	"peanut/internal/service"
	"peanut/internal/template"
	"peanut/internal/validator"
)

func RegisterRegisterHandlers(mux *http.ServeMux, groupService service.GroupService, userService service.UserService) {
	getRegisterHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		templateCtx := templatecontext.GetStandardTemplateContext(r)
		theTemplate := template.GetTemplate("_register")
		err := theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing template:", err)
		}
	})
	mux.Handle("GET /register", getRegisterHandler)

	postRegisterHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username := r.PostFormValue("username")
		email := r.PostFormValue("email")
		password := r.PostFormValue("password")
		password2 := r.PostFormValue("password2")

		if err := validator.ValidateUsername(username); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "Username is not valid: "+err.Error(), w, r)
			return
		}
		if err := validator.ValidateEmail(email); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "Email is not valid: "+err.Error(), w, r)
			return
		}
		if password != password2 {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "Passwords must match.", w, r)
			return
		}
		if err := validator.ValidatePassword(password); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", "Password is not valid: "+err.Error(), w, r)
			return
		}

		userId, err := userService.CreateUser(r, username, email, password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to crete user.", w, r)
			return
		}
		err = groupService.EnrollUserInGroup(r, userId, permgroups.User)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Failed to add new user to default group.", w, r)
			return
		}

		err = requtil.CommitTransactionForRequest(r)
		if err != nil {
			genericpage.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		genericpage.RenderSimpleMessage("Success", "New user has been successfully registered.", w, r)
		logger.Info(r, "Registered user:", userId)
	})
	mux.Handle("POST /register", postRegisterHandler)
}
