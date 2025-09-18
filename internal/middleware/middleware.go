// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"net/http"
	"slices"
)

type MiddlewareFunc func(next http.Handler) http.Handler

func WrapHandler(baseHandler http.Handler, middleware ...MiddlewareFunc) http.Handler {
	result := baseHandler
	// The first middleware wraps the second, which wraps the third, etc
	// Because of this, we need to build the final handler starting at the inner layer and working outwards
	for _, thisMiddleware := range slices.Backward(middleware) {
		result = thisMiddleware(result)
	}
	return result
}
