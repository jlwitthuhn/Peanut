// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"net/http"
	"peanut/internal/endpoints/genericpage"
	"slices"
)

func CheckPermissions(requiredPerms ...string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userPerms, ok := r.Context().Value("userPerms").([]string)
			if !ok {
				genericpage.RenderErrorHttp403Forbidden(w, r)
				return
			}
			for _, requiredPerm := range requiredPerms {
				if slices.Contains(userPerms, requiredPerm) == false {
					genericpage.RenderErrorHttp403Forbidden(w, r)
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
