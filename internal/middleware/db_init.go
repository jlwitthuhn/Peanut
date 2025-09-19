// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"net/http"

	"peanut/internal/logger"
	"peanut/internal/service/db_service"
)

func DatabaseInitCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := db_service.DoesTableExist("config_int")
		if err != nil {
			logger.Fatal(err)
		}
		next.ServeHTTP(w, r)
	})
}
