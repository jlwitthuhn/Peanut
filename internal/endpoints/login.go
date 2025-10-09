// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package endpoints

import (
	"fmt"
	"net/http"
	"peanut/internal/cookie"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/endpoints/templatecontext"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/logger"
	"peanut/internal/service"
)

func RegisterLoginHandlers(mux *http.ServeMux, sessionService service.SessionService) {
	getLoginHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isAlreadyLoggedIn(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		templateCtx := templatecontext.GetStandardTemplateContext(r)
		ep_util.RenderTemplate("_login", templateCtx, w, r)
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
			ep_util.RenderErrorHttp400BadRequestWithMessage(errMsg, w, r)
			return
		}

		err = ep_util.CommitTransactionForRequest(r)
		if err != nil {
			ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to commit transaction.", w, r)
			return
		}

		sesssionCookie := cookie.CreateSessionCookie(sessionId)
		http.SetCookie(w, &sesssionCookie)
		ep_util.RenderSimpleMessage("Success", "You have logged in.", w, r)
	})
	mux.Handle("POST /login", postLoginHandler)
}

func isAlreadyLoggedIn(r *http.Request) bool {
	return r.Context().Value(contextkeys.LoggedIn) == true
}
