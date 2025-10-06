// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"net/http"
	"peanut/internal/cookie"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/security/perms"
	"peanut/internal/service"
)

func Authentication(groupService service.GroupService, sessionService service.SessionService) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookies := r.CookiesNamed(cookie.SessionCookieName)
			for _, thisCookie := range cookies {
				userId, sessionErr := sessionService.GetLoggedInUserIdBySessionId(r, thisCookie.Value)
				if sessionErr != nil {
					continue
				}
				if userId == "" {
					continue
				}
				groups, groupsErr := groupService.GetGroupsByUserId(r, userId)
				if groupsErr != nil {
					continue
				}
				permissions := perms.GetGranularPermissionsForGroups(groups...)
				ctx = context.WithValue(ctx, contextkeys.LoggedIn, true)
				ctx = context.WithValue(ctx, contextkeys.SessionId, thisCookie.Value)
				ctx = context.WithValue(ctx, contextkeys.UserGroups, groups)
				ctx = context.WithValue(ctx, contextkeys.UserId, userId)
				ctx = context.WithValue(ctx, contextkeys.UserPerms, permissions)
				break
			}
			if ctx.Value(contextkeys.LoggedIn) == nil {
				defaultGroups := []string{"Guest"}
				permissions := perms.GetGranularPermissionsForGroups(defaultGroups...)
				ctx = context.WithValue(ctx, contextkeys.LoggedIn, false)
				ctx = context.WithValue(ctx, contextkeys.UserGroups, defaultGroups)
				ctx = context.WithValue(ctx, contextkeys.UserPerms, permissions)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
