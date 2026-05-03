// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package ep_forums

import (
	"net/http"
	"peanut/internal/endpoints/ep_util"
)

func registerForumHomeHandlers(mux *http.ServeMux) {
	getForumHomeHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ep_util.RenderSimpleMessage("Forum", "This page is under construction.", w, r)
	})
	mux.Handle("GET /forum", getForumHomeHandler)
}
