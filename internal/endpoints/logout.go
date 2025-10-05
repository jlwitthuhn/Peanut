// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package endpoints

import (
	"net/http"
	"peanut/internal/cookie"
	"peanut/internal/endpoints/genericpage"
	"peanut/internal/logger"
	"peanut/internal/service"
)

func RegisterLogoutHandlers(mux *http.ServeMux, sessionService service.SessionService) {
	postLogoutHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionId := r.Context().Value("sessionId")
		sessionIdString, ok := sessionId.(string)
		if ok == false {
			logger.Error(r, "'sessionId' is not a string while logging out.")
			w.WriteHeader(http.StatusInternalServerError)
			genericpage.RenderSimpleMessage("Error", "Unable to log out: failed to read session id.", w, r)
			return
		}
		err := sessionService.DestroySession(r, nil, sessionIdString)
		if err != nil {
			logger.Warn(r, "Failed to delete session, proceeding anyways.")
		}

		emptyCookie := cookie.CreateSessionCookie("There was a session id here. It is gone now.")
		http.SetCookie(w, &emptyCookie)
		genericpage.RenderSimpleMessage("Success", "You have been logged out.", w, r)
	})
	mux.Handle("POST /logout", postLogoutHandler)
}
