// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"net/http"
	"peanut/internal/data/datasource"
	"peanut/internal/endpoints/genericpage"
)

func PostgresTransaction() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tx, err := datasource.PostgresHandle().BeginTx(r.Context(), nil)
			if err != nil {
				genericpage.RenderErrorHttp500InternalServerErrorWithMessage("Failed to create database transaction.", w, r)
				return
			}
			defer tx.Rollback()
			ctx := context.WithValue(r.Context(), "tx", tx)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
