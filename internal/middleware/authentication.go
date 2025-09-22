// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"net/http"
	"peanut/internal/cookie"
	"peanut/internal/service"
)

func Authentication(userService service.UserService) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookies := r.CookiesNamed(cookie.SessionCookieName)
			for _, thisCookie := range cookies {
				userId, err := userService.GetLoggedInUserIdBySession(r, nil, thisCookie.Value)
				if err != nil {
					continue
				}
				ctx = context.WithValue(ctx, "userId", userId)
				ctx = context.WithValue(ctx, "loggedIn", true)
				break
			}
			if ctx.Value("loggedIn") == nil {
				ctx = context.WithValue(ctx, "loggedIn", false)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
