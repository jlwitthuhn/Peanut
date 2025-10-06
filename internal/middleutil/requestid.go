// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package middleutil

import (
	"net/http"
	"peanut/internal/keynames/contextkeys"
)

func RetrieveRequestId(r *http.Request) string {
	if r != nil {
		return r.Context().Value(contextkeys.RequestId).(string)
	} else {
		return "--------"
	}
}
