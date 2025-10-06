// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"peanut/internal/keynames/contextkeys"
)

func RequestId() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, contextkeys.RequestId, generateRandomRequestId())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func generateRandomRequestId() string {
	bytes := make([]byte, 4)
	_, err := rand.Read(bytes)
	if err != nil {
		return "ERROR~~~"
	}
	return hex.EncodeToString(bytes)
}
