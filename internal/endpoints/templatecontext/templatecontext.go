// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package templatecontext

import (
	"fmt"
	"net/http"
	"peanut/internal/keynames/contextkeys"
	"peanut/internal/middleutil"
	"strings"
)

func GetStandardTemplateContext(r *http.Request) map[string]any {
	result := make(map[string]any)
	result["LoggedIn"] = r.Context().Value(contextkeys.LoggedIn)

	permSlice, ok := r.Context().Value(contextkeys.UserPerms).([]string)
	if ok {
		for _, perm := range permSlice {
			fullPerm := "Perm_" + strings.ReplaceAll(perm, "/", "_")
			result[fullPerm] = true
		}
	}

	// Always do this one last
	timeMs := middleutil.RequestTimerFinish(r)
	result["RequestDuration"] = fmt.Sprintf("%.1f", timeMs)
	return result
}
