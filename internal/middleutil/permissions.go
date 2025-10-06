// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleutil

import (
	"net/http"
	"peanut/internal/keynames/contextkeys"
	"slices"
)

func RequestHasPermission(r *http.Request, permission string) bool {
	permissionsAny := r.Context().Value(contextkeys.UserPerms)
	permissions := permissionsAny.([]string)
	return slices.Contains(permissions, permission)
}
