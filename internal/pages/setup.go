// This file is part of Peanut and is licensed under the AGPLv3
// https://www.gnu.org/licenses/agpl-3.0.en.html
// SPDX-License-Identifier: AGPL-3.0-only

package pages

import (
	"net/http"
	"peanut/internal/pages/genericpage"
)

func RegisterSetupHandlers(mux *http.ServeMux) {
	getSetupHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		genericpage.RenderSimpleMessage("TODO", "Page not implemented", w, r)
	})
	mux.Handle("GET /setup", getSetupHandler)
}
