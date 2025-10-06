// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"net/http"
	"peanut/internal/keynames/contextkeys"
	"time"
)

func RequestTimer() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextkeys.RequestTimerBegin, time.Now())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
