// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import "net/http"

type MiddlewareFunc func(next http.Handler) http.Handler

func WrapHandler(baseHandler http.Handler, middleware ...MiddlewareFunc) http.Handler {
	var result http.Handler = baseHandler
	// Walk through the list backwards so the order the handlers run matches the list
	for i := len(middleware) - 1; i >= 0; i-- {
		result = middleware[i](result)
	}
	return result
}
