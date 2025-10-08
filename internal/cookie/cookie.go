// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package cookie

import (
	"net/http"
	"strings"
)

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

var unauthenticatedCsrfCookieName = "unauthenticatedCsrf"
var unauthenticatedCsrfTokenLimit = 5

func CreateUnauthenticatedCsrfCookie(req *http.Request, newToken string) http.Cookie {
	tokenList := GetUnauthenticatedCsrfTokens(req)
	tokenList = append([]string{newToken}, tokenList...)
	outLength := min(len(tokenList), unauthenticatedCsrfTokenLimit)
	tokenList = tokenList[:outLength]
	content := strings.Join(tokenList, "|")
	return http.Cookie{
		Name:     unauthenticatedCsrfCookieName,
		Value:    content,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
}

func GetUnauthenticatedCsrfTokens(req *http.Request) []string {
	cookies := req.CookiesNamed(unauthenticatedCsrfCookieName)
	for _, thisCookie := range cookies {
		return strings.Split(thisCookie.Value, "|")
	}
	return []string{}
}
