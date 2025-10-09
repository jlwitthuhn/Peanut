// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"net/http"
	"peanut/internal/cookie"
	"peanut/internal/endpoints/ep_util"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/keynames/sessionkeys"
	"peanut/internal/security"
	"peanut/internal/service"
	"slices"
)

func CsrfProtection(sessionService service.SessionService) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			loggedIn, ok := r.Context().Value(contextkeys.LoggedIn).(bool)
			if !ok {
				ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to read logged in status.", w, r)
				return
			}

			ctx := r.Context()
			if loggedIn {
				sessionId, ok := r.Context().Value(contextkeys.SessionId).(string)
				if !ok {
					ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Failed to read session id.", w, r)
					return
				}
				token, err := sessionService.GetString(r, sessionId, sessionkeys.CsrfToken)
				if err != nil {
					ep_util.RenderErrorHttp500InternalServerErrorWithMessage("Database does not contain a valid CSRF token for this session.", w, r)
					return
				}
				ctx = context.WithValue(ctx, contextkeys.CsrfToken, token)

				if r.Method == http.MethodPost {
					submittedToken := r.PostFormValue("csrf")
					if submittedToken != token {
						ep_util.RenderErrorHttp400BadRequestWithMessage("Invalid CSRF token.", w, r)
						return
					}
				}
			} else {
				existingTokens := cookie.GetUnauthenticatedCsrfTokens(r)
				if len(existingTokens) == 1 && len(existingTokens[0]) == 0 {
					ep_util.RenderErrorHttp400BadRequestWithMessage("No valid CSRF token in cookies.", w, r)
					return
				}

				if r.Method == http.MethodPost {
					submittedToken := r.PostFormValue("csrf")
					if !slices.Contains(existingTokens, submittedToken) {
						ep_util.RenderErrorHttp400BadRequestWithMessage("Invalid CSRF token.", w, r)
						return
					}
				}

				// Do this last so we don't set a new token on failed POSTs
				newToken := security.GenerateCsrfSmallToken()
				csrfCookie := cookie.CreateUnauthenticatedCsrfCookie(r, newToken)
				http.SetCookie(w, &csrfCookie)

				ctx = context.WithValue(ctx, contextkeys.CsrfToken, newToken)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
