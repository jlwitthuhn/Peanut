// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"net/http"
	"peanut/internal/endpoints/genericpage"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/keynames/sessionkeys"
	"peanut/internal/service"
)

func CsrfProtection(sessionService service.SessionService) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggedIn, ok := r.Context().Value(contextkeys.LoggedIn).(bool)
			if !ok || !loggedIn {
				next.ServeHTTP(w, r)
				return
			}

			sessionId, ok := r.Context().Value(contextkeys.SessionId).(string)
			if !ok {
				genericpage.RenderErrorHttp500InternalServerErrorWithMessage("Failed to read session id.", w, r)
				return
			}
			token, err := sessionService.GetString(r, sessionId, sessionkeys.CsrfToken)
			if err != nil {
				genericpage.RenderErrorHttp500InternalServerErrorWithMessage("Database does not contain a valid CSRF token for this session.", w, r)
				return
			}
			ctx := r.Context()
			ctx = context.WithValue(ctx, contextkeys.CsrfToken, token)

			if r.Method == http.MethodPost {
				submittedToken := r.PostFormValue("csrf")
				if submittedToken != token {
					genericpage.RenderErrorHttp400BadRequestWithMessage("Invalid CSRF token.", w, r)
					return
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
