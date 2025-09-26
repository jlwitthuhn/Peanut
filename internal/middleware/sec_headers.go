// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import "net/http"

func SecurityHeaders() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Don't allow embedding
			w.Header().Set("Content-Security-Policy", "frame-ancestors 'none';")
			w.Header().Set("X-Frame-Options", "DENY")

			next.ServeHTTP(w, r)
		})
	}
}
