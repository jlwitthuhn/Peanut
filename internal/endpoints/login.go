// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package endpoints

import (
	"fmt"
	"net/http"
	"peanut/internal/cookie"
	"peanut/internal/endpoints/genericpage"
	"peanut/internal/endpoints/requtil"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/logger"
	"peanut/internal/service"
	"peanut/internal/template"
)

func RegisterLoginHandlers(mux *http.ServeMux, sessionService service.SessionService) {
	getLoginHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAlreadyLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		theTemplate := template.GetTemplate("_login")
		err := theTemplate.Execute(w, templateCtx)
		if err != nil {
			logger.Error(r, "Error executing template:", err)
		}
	})
	mux.Handle("GET /login", getLoginHandler)

	postLoginHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAlreadyLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		username := r.PostFormValue("username")
		password := r.PostFormValue("password")
		sessionId, err := sessionService.CreateSession(r, username, password)
		if err != nil {
			logger.Error(r, "Error creating session:", err)
			errMsg := fmt.Sprint("Error logging in: ", err)
			w.WriteHeader(http.StatusBadRequest)
			genericpage.RenderSimpleMessage("Error", errMsg, w, r)
			return
		}

		err = requtil.CommitTransactionForRequest(r)
		if err != nil {
			genericpage.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		sesssionCookie := cookie.CreateSessionCookie(sessionId)
		http.SetCookie(w, &sesssionCookie)
		genericpage.RenderSimpleMessage("Success", "You have logged in.", w, r)
	})
	mux.Handle("POST /login", postLoginHandler)
}

func isAlreadyLoggedIn(r *http.Request) bool {
	return r.Context().Value("loggedIn") == true
}
