// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package templatecontext

import (
	"net/http"
	"peanut/internal/middleutil"
)

func GetStandardTemplateContext(r *http.Request) map[string]any {
	result := make(map[string]any)
	result["LoggedIn"] = r.Context().Value("loggedIn")

	// Always do this one last
	result["RequestDuration"] = middleutil.RequestTimerFinish(r)

	return result
}
