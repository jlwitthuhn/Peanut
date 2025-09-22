// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package cookie

import "net/http"

var SessionCookieName = "session"

func CreateSessionCookie(sessionId string) http.Cookie {
	return http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionId,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}
