// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleware

import (
	"context"
	"net/http"
	"peanut/internal/cookie"
	"peanut/internal/data/datasource"
	"peanut/internal/pages/genericpage"
	"peanut/internal/security/perms"
	"peanut/internal/service"
)

func Authentication(groupService service.GroupService, userService service.UserService) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			cookies := r.CookiesNamed(cookie.SessionCookieName)
			tx, txErr := datasource.PostgresHandle().BeginTx(r.Context(), nil)
			if txErr != nil {
				w.WriteHeader(http.StatusInternalServerError)
				genericpage.RenderSimpleMessage("Authentication Error", "Failed to create database transaction.", w, r)
				return
			}
			defer tx.Rollback()
			for _, thisCookie := range cookies {
				userId, sessionErr := userService.GetLoggedInUserIdBySession(r, nil, thisCookie.Value)
				if sessionErr != nil {
					continue
				}
				if userId == "" {
					continue
				}
				groups, groupsErr := groupService.GetGroupsByUserId(tx, userId)
				if groupsErr != nil {
					continue
				}
				permissions := perms.GetGranularPermissionsForGroups(groups...)
				ctx = context.WithValue(ctx, "loggedIn", true)
				ctx = context.WithValue(ctx, "sessionId", thisCookie.Value)
				ctx = context.WithValue(ctx, "userGroups", groups)
				ctx = context.WithValue(ctx, "userId", userId)
				ctx = context.WithValue(ctx, "userPerms", permissions)
				break
			}
			if ctx.Value("loggedIn") == nil {
				defaultGroups := []string{"Guest"}
				permissions := perms.GetGranularPermissionsForGroups(defaultGroups...)
				ctx = context.WithValue(ctx, "loggedIn", false)
				ctx = context.WithValue(ctx, "userGroups", defaultGroups)
				ctx = context.WithValue(ctx, "userPerms", permissions)
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
